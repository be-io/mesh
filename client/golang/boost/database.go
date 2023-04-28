/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"io"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// --dsn=user:password@tcp(127.0.0.1:3306)/dbname  // root:@tcp(127.0.0.1:3306)/mesh
//go:generate go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
//go:generate sqlc generate -f mysql/sqlc.yaml -x
//go:generate sqlc generate -f postgresql/sqlc.yaml -x

var (
	sources = &databases{dbs: map[string]*database{}}
	pattern = regexp.MustCompile(`tcp\(([\w|\W]+)\)`)
)

type database struct {
	*sql.DB
	schema string
}

type databases struct {
	dbs map[string]*database
	sync.RWMutex
}

func (that *databases) decryptPwd(ctx context.Context, uri *url.URL) (string, error) {
	passwd, _ := uri.User.Password()
	if "" == passwd || !strings.Contains(passwd, "ENC") {
		return passwd, nil
	}
	mtx := mpc.ContextWith(ctx).Resume(ctx)
	mtx.SetAttribute(mpc.AddressKey, macro.Env("gaia-janus", "JSN", "jsn"))
	mtx.SetAttribute(mpc.TimeoutKey, time.Second*3)
	mtx.GetPrincipals().Push(&types.Principal{NodeId: types.LocalNodeId, InstId: types.LocalInstId})
	defer func() { mtx.GetPrincipals().Pop() }()
	if strings.Index(passwd, "ENC(") == 0 && strings.Index(passwd, ")") == len(passwd)-1 {
		passwd = passwd[4 : len(passwd)-1]
	}
	ret, err := aware.Dispatcher.InvokeLRG(mtx, types.LocURN(ctx, "janus.open.invoke"), map[string]map[string]interface{}{
		"input": {
			"method":  "com.trustbe.janus.icbc.decrypt",
			"content": fmt.Sprintf("{\"input\":\"%s\"}", passwd),
		}})
	if nil != err {
		return passwd, cause.Error(err)
	}
	str, err := aware.JSON.EncodeString(ret)
	if nil != err {
		return passwd, cause.Error(err)
	}
	log.Info(ctx, str)
	var output map[string]string
	if _, err = aware.JSON.DecodeString(str, &output); nil != err {
		return passwd, cause.Error(err)
	}
	if output["code"] != "200" {
		return passwd, cause.Errorf("Decrypt password failed. ")
	}
	var content map[string]string
	if _, err = aware.JSON.DecodeString(output["content"], &content); nil != err {
		return passwd, cause.Error(err)
	}
	return tool.Anyone(content["output"], passwd), nil
}

func (that *databases) open(ctx context.Context, dsn string) (*database, error) {
	if !strings.Contains(dsn, "://") {
		return nil, cause.Errorf("DSN must like mysql://username:password@hostname:port/database?parseTime=true")
	}
	if db := func() *database {
		that.RLock()
		defer that.RUnlock()
		return that.dbs[dsn]
	}(); nil != db {
		return db, nil
	}
	that.Lock()
	defer that.Unlock()
	if nil != that.dbs[dsn] {
		return that.dbs[dsn], nil
	}
	uri, err := url.Parse(dsn)
	if nil != err {
		return nil, cause.Error(err)
	}
	if "true" != uri.Query().Get("parseTime") && strings.Contains(dsn, "mysql") {
		qr := uri.Query()
		qr.Set("parseTime", "true")
		uri.RawQuery = qr.Encode()
	}
	if "" == uri.Query().Get("loc") && strings.Contains(dsn, "mysql") {
		qr := uri.Query()
		qr.Set("loc", time.Now().Local().Location().String())
		uri.RawQuery = qr.Encode()
	}
	db, err := func() (*sql.DB, error) {
		passwd, err := that.decryptPwd(ctx, uri)
		if nil != err {
			return nil, cause.Error(err)
		}
		if strings.Contains(uri.Scheme, "postgres") {
			uri.User = url.UserPassword(uri.User.Username(), passwd)
			log.Info(ctx, "DSN=%s", uri.String())
			return sql.Open("postgres", uri.String())
		}
		if strings.Contains(uri.Scheme, "mysql") {
			adds := fmt.Sprintf("%s:%s@tcp(%s)%s?%s", uri.User.Username(), passwd, uri.Host, uri.Path, uri.RawQuery)
			log.Info(ctx, "DSN=%s", adds)
			return sql.Open("mysql", adds)
		}
		uri.User = url.UserPassword(uri.User.Username(), passwd)
		log.Info(ctx, "DSN=%s", uri.String())
		return sql.Open(uri.Scheme, uri.String())
	}()
	if nil != err {
		return nil, cause.Error(err)
	}
	db.SetConnMaxLifetime(time.Second * 12)
	db.SetConnMaxIdleTime(time.Second * 12)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(30)

	that.dbs[dsn] = &database{DB: db, schema: uri.Scheme}
	return that.dbs[dsn], nil
}

// Conn open the connection.
// Open &sql.TxOptions{Isolation: sql.LevelReadCommitted}
func (that *databases) Conn(ctx context.Context, dsn string) (string, *sql.Conn, error) {
	if "" == dsn {
		return "", nil, cause.Errorf("DSN must set in environment.")
	}
	db, err := that.open(ctx, dsn)
	if nil != err {
		return "", nil, cause.Error(err)
	}
	mtx, cancel := context.WithTimeout(macro.Context(), time.Second*10)
	defer cancel()
	conn, err := db.Conn(mtx)
	return db.schema, conn, err
}

type Session[T io.Closer] struct {
	dsn  func() string
	next func(ctx context.Context, schema string, conn *sql.Conn, tx *sql.Tx) (T, error)
}

func (that *Session[T]) NT(ctx context.Context, closure func(session T) error) error {
	_, err := RT(ctx, that, func(session T) (any, error) {
		return nil, closure(session)
	})
	return cause.Error(err)
}

func RT[T io.Closer, V any](ctx context.Context, session *Session[T], closure func(session T) (V, error)) (v V, err error) {
	schema, conn, err := sources.Conn(ctx, session.dsn())
	if nil != err {
		log.Warn(ctx, "Connect %s, %s ", session.dsn(), err.Error())
		return v, cause.Errorf("Connect database err: %s", err.Error())
	}
	defer func() {
		log.Catch(conn.Close())
	}()
	nx, err := session.next(ctx, schema, conn, nil)
	if nil != err {
		return v, cause.Error(err)
	}
	defer func() { log.Catch(nx.Close()) }()
	return closure(nx)
}

func (that *Session[T]) NTX(ctx context.Context, closure func(session T) error) error {
	_, err := RTX(ctx, that, func(session T) (any, error) {
		return nil, closure(session)
	})
	return cause.Error(err)
}

func RTX[T io.Closer, V any](ctx context.Context, session *Session[T], closure func(session T) (V, error)) (v V, err error) {
	schema, conn, err := sources.Conn(ctx, session.dsn())
	if nil != err {
		log.Warn(ctx, "Connect %s, %s ", session.dsn(), err.Error())
		return v, cause.Errorf("Connect database err: %s", err.Error())
	}
	defer func() {
		log.Catch(conn.Close())
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if nil != err {
		return v, cause.Error(err)
	}
	defer func() {
		if cover := recover(); nil != cover {
			log.Catch(tx.Rollback())
			log.Error(ctx, "Commit transaction unexpected, %v", cover)
			log.Error(ctx, string(debug.Stack()))
			err = cause.Errorf("%v", cover)
		}
	}()

	accessor, err := session.next(ctx, schema, conn, tx)
	if nil != err {
		return v, cause.Error(err)
	}
	defer func() { log.Catch(accessor.Close()) }()

	v, err = closure(accessor)
	if nil != err {
		log.Catch(tx.Rollback())
	} else {
		log.Catch(tx.Commit())
	}
	return v, cause.Error(err)
}

// DAL construct sql session, tx maybe nil if no transaction required.
func DAL[T io.Closer](dsn func() string, next func(ctx context.Context, schema string, conn *sql.Conn, tx *sql.Tx) (T, error)) *Session[T] {
	return &Session[T]{dsn: func() string {
		x := dsn()
		if strings.Contains(x, "://") {
			return x
		}
		return fmt.Sprintf("mysql://%s", pattern.ReplaceAllString(x, "$1"))
	}, next: next}
}

/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"context"
	"database/sql"
	_ "github.com/dolthub/go-mysql-server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kyleconroy/sqlc"
	_ "github.com/lib/pq"
	"github.com/opendatav/mesh/client/golang/boost"
	"github.com/opendatav/mesh/plugin/metabase/dal"
	"strings"
)

// --dsn=user:password@tcp(127.0.0.1:3306)/dbname  // root:@tcp(127.0.0.1:3306)/mesh
//go:generate go install github.com/ducesoft/sqlc/cmd/sqlc@latest
//go:generate sqlc generate -f mysql/sqlc.yaml -x
//go:generate sqlc generate -f postgresql/sqlc.yaml -x

var sessions = boost.DAL(func() string { return metabase.DSN }, func(ctx context.Context, schema string, conn *sql.Conn, tx *sql.Tx) (dal.DAL, error) {
	if nil != tx {
		if strings.Contains(schema, "postgres") {
			return dal.NewPostgresql(tx).WithTx(tx), nil
		}
		return dal.NewMysql(tx).WithTx(tx), nil
	}
	if strings.Contains(schema, "postgres") {
		return dal.NewPostgresql(conn), nil
	}
	return dal.PrepareMysql(ctx, conn)
})

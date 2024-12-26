/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"fmt"
	"github.com/opendatav/mesh/client/golang/log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Vars map[string]string

func (that Vars) GetDuration(opt string) time.Duration {
	val := that.GetString(opt)
	if val == "" {
		return time.Duration(0)
	}
	if strings.Contains(val, "d") {
		val = strings.Replace(val, "d", "", 1)
		days, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return time.Duration(0)
		}
		return time.Hour * 24 * time.Duration(days)
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return time.Duration(0)
	}
	return d
}

func (that Vars) GetBool(opt string) bool {
	val := that.GetString(opt)
	if val == "" {
		return false
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		log.Warn0("Unable to parse %s as bool for key: %s. Options: %s\n", val, opt)
	}
	return b
}

func (that Vars) GetFloat64(opt string) float64 {
	val := that.GetString(opt)
	if val == "" {
		return 0
	}
	floatV, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Warn0("Unable to parse %s as float64 for key: %s. Options: %s\n", val, opt)
	}
	return floatV
}

func (that Vars) GetInt64(opt string) int64 {
	val := that.GetString(opt)
	if val == "" {
		return 0
	}
	i, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		log.Warn0("Unable to parse %s as int64 for key: %s. Options: %s\n", val, opt)
	}
	return i
}

func (that Vars) GetUint64(opt string) uint64 {
	val := that.GetString(opt)
	if val == "" {
		return 0
	}
	u, err := strconv.ParseUint(val, 0, 64)
	if err != nil {
		log.Warn0("Unable to parse %s as uint64 for key: %s. Options: %s\n", val, opt)
	}
	return u
}

func (that Vars) GetUint32(opt string) uint32 {
	val := that.GetString(opt)
	if val == "" {
		return 0
	}
	u, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		log.Warn0("Unable to parse %s as uint32 for key: %s. Options: %s\n, %s", val, opt, err.Error())
	}
	return uint32(u)
}

func (that Vars) GetString(opt string) string {
	if that == nil {
		return ""
	}
	return that[opt]
}

func (that Vars) GetPath(opt string) string {
	p := that.GetString(opt)
	path, err := expandPath(p)
	if err != nil {
		log.Error0("Failed to get path: %+v", err)
	}
	return path
}

// expandPath expands the paths containing ~ to /home/user. It also computes the absolute path
// from the relative paths. For example: ~/abc/../cef will be transformed to /home/user/cef.
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return "", nil
	}
	if path[0] == '~' && (len(path) == 1 || os.IsPathSeparator(path[1])) {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("Failed to get the home directory of the user, %s ", err.Error())
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}

	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Failed to generate absolute path, %s ", err.Error())
	}
	return path, nil
}

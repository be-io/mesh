/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import "strings"

const (
	CommitID = "commit_id"
	OS       = "os"
	ARCH     = "arch"
	Version  = "version"
)

type Versions struct {
	Version string            `index:"0" json:"version,omitempty" yaml:"version,omitempty" xml:"version,omitempty" comment:"Platform version. main.sub.feature.bugfix like 1.5.0.0"`
	Infos   map[string]string `index:"5" json:"infos,omitempty" yaml:"infos,omitempty" xml:"infos,omitempty" comment:"Network modules info."`
}

func (that *Versions) Compare(v string) int {
	if "" == that.Version {
		return -1
	}
	if "" == v {
		return 1
	}
	rvs := strings.Split(that.Version, ".")
	vs := strings.Split(v, ".")
	for index := 0; index < len(rvs); index++ {
		if len(vs) <= index {
			return 0
		}
		if "*" == vs[index] {
			continue
		}
		c := strings.Compare(rvs[index], vs[index])
		if 0 != c {
			return c
		}
	}
	return 0
}

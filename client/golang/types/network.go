/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type NetWork struct {
	ID            string `json:"id" yaml:"id"`
	Name          string `json:"name" yaml:"name"`
	Description   string `json:"description" yaml:"description"`
	Certification string `json:"certification" yaml:"certification"`
	IP            string `json:"ip" yaml:"ip"`
	Port          string `json:"port" yaml:"port"`
	CreateAt      int64  `json:"create_at" yaml:"create_at"`
	UpdateAt      int64  `json:"update_at" yaml:"update_at"`
}

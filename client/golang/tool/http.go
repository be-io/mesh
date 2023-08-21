/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     time.Minute * 10,
	},
	Timeout: time.Minute * 2,
}

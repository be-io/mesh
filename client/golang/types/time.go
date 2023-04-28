/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"encoding/json"
	"github.com/be-io/mesh/client/golang/cause"
	"strconv"
	"time"
)

func init() {
	var _ json.Marshaler = new(Time)
	var _ json.Unmarshaler = new(Time)
	var _ json.Marshaler = new(Duration)
	var _ json.Unmarshaler = new(Duration)
}

type Time time.Time

func (that *Time) UnmarshalJSON(bytes []byte) error {
	timestamp, err := strconv.ParseInt(string(bytes), 10, 64)
	if nil != err {
		return cause.Error(err)
	}
	*that = Time(time.UnixMilli(timestamp))
	return nil
}

func (that Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(that).UnixMilli(), 10)), nil
}

type Duration time.Duration

func (that *Duration) UnmarshalJSON(bytes []byte) error {
	timestamp, err := strconv.ParseInt(string(bytes), 10, 64)
	if nil != err {
		return cause.Error(err)
	}
	*that = Duration(timestamp * time.Millisecond.Nanoseconds())
	return nil
}

func (that Duration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Duration(that).Milliseconds(), 10)), nil
}

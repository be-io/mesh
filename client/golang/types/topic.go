/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import "github.com/opendatav/mesh/client/golang/macro"

type Topic struct {

	// Topic is message topic
	Topic string `index:"0" json:"topic" xml:"topic" yaml:"topic"`
	// Code is message topic
	Code string `index:"5" json:"code" xml:"code" yaml:"code"`
	// Group is message group
	Group string `index:"10" json:"group" xml:"group" yaml:"group"`
	// Group is message group
	Sets string `index:"15" json:"sets" xml:"sets" yaml:"sets"`
}

func (that *Topic) Match(bindings ...*macro.Btt) bool {
	for _, binding := range bindings {
		if nil == binding {
			continue
		}
		if that.Topic == binding.Topic && (that.Code == binding.Code || binding.Code == "*") {
			return true
		}
	}
	return false
}

func (that *Topic) With(binding *macro.Btt) *Topic {
	that.Topic = binding.Topic
	that.Code = binding.Code
	return that
}

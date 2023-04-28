/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	MessageVersion = "1.0.0"
)

type Event struct {
	Version   string     `index:"0" json:"version" yaml:"version" xml:"version" comment:""`
	Tid       string     `index:"5" json:"tid" yaml:"tid" xml:"tid" comment:""`
	Sid       string     `index:"10" json:"sid" yaml:"sid" xml:"sid" comment:""`
	Eid       string     `index:"15" json:"eid" yaml:"eid" xml:"eid" comment:""`
	Mid       string     `index:"20" json:"mid" yaml:"mid" xml:"mid" comment:""`
	Timestamp string     `index:"25" json:"timestamp" yaml:"timestamp" xml:"timestamp" comment:""`
	Source    *Principal `index:"30" json:"source" yaml:"source" xml:"source" comment:""`
	Target    *Principal `index:"35" json:"target" yaml:"target" xml:"target" comment:""`
	Binding   *Topic     `index:"40" json:"binding" yaml:"binding" xml:"binding" comment:""`
	Entity    *Entity    `index:"45" json:"entity" xml:"entity" yaml:"entity" comment:""`
}

func (that *Event) GetObject() (interface{}, error) {
	if nil == that.Entity {
		return nil, nil
	}
	return that.Entity.ReadObject()
}

func (that *Event) TryGetObject(ptr interface{}) error {
	if nil == that.Entity {
		return nil
	}
	return that.Entity.TryReadObject(ptr)
}

/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Document struct {
	Metadata  map[string]string `index:"0" json:"metadata" xml:"metadata" yaml:"metadata"`
	Content   string            `index:"5" json:"content" xml:"content" yaml:"content"`
	Timestamp int64             `index:"10" json:"timestamp" xml:"timestamp" yaml:"timestamp"`
}

type DocumentMetadata struct {
	Queries   map[string]string `index:"0" json:"queries" xml:"queries" yaml:"queries"`
	Start     float64           `index:"5" json:"start" xml:"start" yaml:"start"`
	End       float64           `index:"10" json:"end" xml:"end" yaml:"end"`
	Limit     int64             `index:"15" json:"limit" xml:"limit" yaml:"limit"`
	Step      string            `index:"20" json:"step" xml:"step" yaml:"step"`
	Direction string            `index:"25" json:"direction" xml:"direction" yaml:"direction"`
	RouteKey  string            `index:"30" json:"route_key" xml:"route_key" yaml:"route_key"` // across node read
}

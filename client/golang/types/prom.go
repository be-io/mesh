/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type PromRangeQuery struct {
	IndexName string  `json:"indexName"`
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	Step      int     `json:"step"`
	PartyId   string  `json:"partyId"`
}

type TasksIndex struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	FTaskType     string `json:"fTaskType"`
	FPartyId      string `json:"FPartyId"`
	Initiator     string `json:"initiator"`
	JobId         string `json:"jobId"`
	StartTime     int64  `json:"startTime"`
	EndTime       int64  `json:"endTime"`
	Sort          string `json:"sort"`
	PageNo        int    `json:"pageNo"`
	PageSize      int    `json:"pageSize"`
	Authorization string `json:"authorization"`
}

type LogsIndex struct {
	ComponentName string `json:"componentName,omitempty"`
	JobId         string `json:"jobId,omitempty"`
	PartyId       string `json:"partyId,omitempty"`
	Role          string `json:"role,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Start         int64  `json:"start"`
	End           int64  `json:"end"`
	Step          int    `json:"step"`
	Limit         int64  `json:"limit"`
	// Supported values are forward or backward. Defaults to backward.
	Direction string `json:"direction,omitempty"`
	RetType   string `json:"retType,omitempty"`
	LogSource string `json:"logSource,omitempty"`
	Page      int64  `json:"page"`
	Size      int64  `json:"size"`
}

type PromProxy struct {
	// Path prom GET API path
	Path string `json:"path"`
	// Params prom GET API query params
	Params map[string]interface{} `json:"params"`
}

type PromLabelValue struct {
	Label string `json:"label"`
}

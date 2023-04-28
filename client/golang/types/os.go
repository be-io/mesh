/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type OSCharts struct {
	ReleaseName  string
	RepoName     string
	ChartName    string
	K8sNamespace string
	ChartArgs    map[string]string
	Cfs          []byte
	ChartType    string
	UniqueName   string
	KubeConfig   []byte
}

/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

var (
	WorkStart  VertexKind = 1
	WorkFinish VertexKind = 2
	WorkJob    VertexKind = 4
	WorkMan    VertexKind = 8
	WorkWay    VertexKind = 16
	WorkTimer  VertexKind = 32
)

type VertexKind int64

type WorkIntent struct {
	BNO     string            `index:"0" json:"bno" yaml:"bno" xml:"bno" comment:"Business code"`
	CNO     string            `index:"1" json:"cno" yaml:"cno" xml:"cno" comment:"Workflow chart code"`
	Context map[string]string `index:"2" json:"context" yaml:"context" xml:"context" comment:"Workflow context"`
	Applier *Worker           `index:"3" json:"applier" yaml:"applier" xml:"applier" comment:"Workflow applier"`
}

type WorkVertex struct {
	Name  string            `index:"0" json:"name" yaml:"name" xml:"name" comment:"Workflow name"`
	Alias string            `index:"1" json:"alias" yaml:"alias" xml:"alias" comment:"Workflow alias"`
	Attrs map[string]string `index:"2" json:"attrs" yaml:"attrs" xml:"attrs" comment:"Workflow vertex attributes"`
	Kind  int64             `index:"3" json:"kind" yaml:"kind" xml:"kind" comment:"Workflow vertex kind"`
	Group string            `index:"4" json:"group" yaml:"group" xml:"group" comment:"Workflow review group code"`
}

type WorkSide struct {
	Src       string `index:"0" json:"src" yaml:"src" xml:"src" comment:"Workflow side src name"`
	Dst       string `index:"1" json:"dst" yaml:"dst" xml:"dst" comment:"Workflow side dst name"`
	Condition string `index:"2" json:"condition" yaml:"condition" xml:"condition" comment:"Workflow side condition"`
}

type WorkChart struct {
	CNO        string        `index:"0" json:"cno" yaml:"cno" xml:"cno" comment:"Workflow chart code"`
	Name       string        `index:"1" json:"name" yaml:"name" xml:"name" comment:"Workflow name"`
	Vertices   []*WorkVertex `index:"2" json:"vertices" yaml:"vertices" xml:"vertices" comment:"Workflow vertices"`
	Sides      []*WorkSide   `index:"3" json:"sides" yaml:"sides" xml:"sides" comment:"Workflow sides"`
	Status     int64         `index:"4" json:"status" yaml:"status" xml:"status" comment:"Workflow status"`
	Maintainer *Worker       `index:"5" json:"maintainer" yaml:"maintainer" xml:"maintainer" comment:"Workflow maintainer"`
}

func (that *WorkChart) Vertex(vertex *WorkVertex) {
	that.Vertices = append(that.Vertices, vertex)
}

func (that *WorkChart) Link(side *WorkSide) {
	that.Sides = append(that.Sides, side)
}

type WorkRoutine struct {
	RNO     string            `index:"0" json:"rno" yaml:"rno" xml:"rno" comment:"Workflow routine code"`
	BNO     string            `index:"1" json:"bno" yaml:"bno" xml:"bno" comment:"Business code"`
	Context map[string]string `index:"2" json:"context" yaml:"context" xml:"context" comment:"Workflow context"`
	Status  int64             `index:"3" json:"status" yaml:"status" xml:"status" comment:"Workflow status"`
	Chart   *WorkChart        `index:"4" json:"chart" yaml:"chart" xml:"chart" comment:"Workflow chart"`
	Tasks   []*WorkTask       `index:"5" json:"tasks" yaml:"tasks" xml:"tasks" comment:"Workflow tasks"`
}

type WorkTask struct {
	Vertex    *WorkVertex       `index:"0" json:"vertex" yaml:"vertex" xml:"vertex" comment:"Workflow vertex"`
	Reviewers []*Worker         `index:"1" json:"reviewers" yaml:"reviewers" xml:"reviewers" comment:"Workflow vertex reviewers"`
	Status    int64             `index:"2" json:"status" yaml:"status" xml:"status" comment:"Workflow vertex status"`
	Context   map[string]string `index:"3" json:"context" yaml:"context" xml:"context" comment:"Workflow context"`
}

type WorkGroup struct {
	NO      string    `index:"0" json:"no" yaml:"no" xml:"no" comment:"Work group identity"`
	Name    string    `index:"1" json:"name" yaml:"name" xml:"name" comment:"Work group name"`
	Status  int64     `index:"2" json:"status" yaml:"status" xml:"status" comment:"Workflow group status"`
	Workers []*Worker `index:"3" json:"workers" yaml:"workers" xml:"workers" comment:"Work group workers"`
}

type Worker struct {
	NO    string `index:"0" json:"no" yaml:"no" xml:"no" comment:"Worker identity"`
	Name  string `index:"1" json:"name" yaml:"name" xml:"name" comment:"Worker name"`
	Alias string `index:"2" json:"alias" yaml:"alias" xml:"alias" comment:"Worker alias"`
}

type WorkContext struct {
	RNO      string            `index:"0" json:"rno" yaml:"rno" xml:"rno" comment:"Workflow routine code"`
	BNO      string            `index:"1" json:"bno" yaml:"bno" xml:"bno" comment:"Business code"`
	CNO      string            `index:"2" json:"cno" yaml:"cno" xml:"cno" comment:"Workflow chart code"`
	Context  map[string]string `index:"3" json:"context" yaml:"context" xml:"context" comment:"Workflow context"`
	Vertex   *WorkVertex       `index:"4" json:"vertex" yaml:"vertex" xml:"vertex" comment:"Workflow vertex"`
	Task     WorkTask          `index:"5" json:"task" yaml:"task" xml:"task" comment:"Workflow task"`
	Applier  *Worker           `index:"6" json:"applier" yaml:"applier" xml:"applier" comment:"Workflow applier"`
	Reviewer *Worker           `index:"7" json:"reviewer" yaml:"reviewer" xml:"reviewer" comment:"Workflow reviewer"`
}

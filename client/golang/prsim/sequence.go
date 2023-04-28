/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/tool"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ISequence = (*Sequence)(nil)

// Sequence spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Sequence interface {

	// Next
	// @MPI("mesh.sequence.next")
	Next(ctx context.Context, kind string, length int) (string, error)

	// Section
	// @MPI("mesh.sequence.section")
	Section(ctx context.Context, kind string, size int, length int) ([]string, error)
}

const (
	DefaultSectionSize    = 5
	DefaultSequenceLength = 8
	DefaultRetries        = 5
)

type SectionSyncer[T any] interface {
	Tx(ctx context.Context, tx func(session T) ([]string, error)) ([]string, error)
	Incr(ctx context.Context, session T, kind string) (*Section, error)
	Init(ctx context.Context, session T, section *Section) error
	Sync(ctx context.Context, session T, section *Section) error
}

type Section struct {
	Kind    string `json:"kind"`
	Min     int64  `json:"min"`
	Max     int64  `json:"max"`
	Size    int32  `json:"size"`
	Length  int32  `json:"length"`
	Version int32  `json:"version"`
}

type SyncSequence[T any] struct {
	sync.RWMutex
	Sections map[string]chan string
	Syncer   SectionSyncer[T]
	Macro    *macro.Att
}

func (that *SyncSequence[T]) Att() *macro.Att {
	return that.Macro
}

func (that *SyncSequence[T]) TryNextSection(ctx context.Context, kind string) chan string {
	if section := func() chan string {
		that.RLock()
		defer that.RUnlock()
		return that.Sections[kind]

	}(); nil != section {
		return section
	}
	that.Lock()
	defer that.Unlock()
	if nil != that.Sections[kind] {
		return that.Sections[kind]
	}
	that.Sections[kind] = make(chan string, 31)
	return that.Sections[kind]
}

func (that *SyncSequence[T]) Next(ctx context.Context, kind string, length int) (string, error) {
	section := that.TryNextSection(ctx, kind)
	for index := 0; ; index++ {
		select {
		case nv, ok := <-section:
			if !ok {
				return "", cause.Errorf("Sequence pool has been closed. ")
			}
			return nv, nil
		case <-time.After(time.Millisecond * time.Duration(index*30)):
			sections, err := that.TrySection(ctx, kind, length)
			if nil != err && index > DefaultRetries {
				return "", cause.Error(err)
			}
			if nil != err {
				continue
			}
			for _, sec := range sections {
				section <- sec
			}
		}
	}
}

func (that *SyncSequence[T]) Section(ctx context.Context, kind string, size int, length int) ([]string, error) {
	for index := 0; ; index++ {
		sections, err := that.TrySection(ctx, kind, length)
		if nil != err && index > DefaultRetries {
			return nil, cause.Error(err)
		}
		if nil != err {
			continue
		}
		return sections, nil
	}
}

func (that *SyncSequence[T]) TrySection(ctx context.Context, kind string, length int) ([]string, error) {
	return that.Syncer.Tx(ctx, func(session T) ([]string, error) {
		var sections []string
		section, err := that.Syncer.Incr(ctx, session, kind)
		// sql.ErrNoRows
		if nil != err && sql.ErrNoRows != err {
			return nil, cause.Error(err)
		}
		if nil == section || sql.ErrNoRows == err {
			xl := int32(tool.Ternary(length < 1, DefaultSequenceLength, length))
			section = &Section{
				Kind:    kind,
				Min:     DefaultSectionSize,
				Max:     math.MaxInt,
				Size:    DefaultSectionSize,
				Length:  xl,
				Version: 0,
			}
			if err = that.Syncer.Init(ctx, session, section); nil != err {
				return nil, cause.Error(err)
			}
			section.Size = DefaultSectionSize
			section.Length = xl
			section.Min = 0
		} else {
			if err = that.Syncer.Sync(ctx, session, &Section{
				Kind:    section.Kind,
				Min:     section.Min + int64(section.Size),
				Max:     section.Max,
				Size:    section.Size,
				Length:  section.Length,
				Version: section.Version,
			}); nil != err {
				return nil, cause.Error(err)
			}
		}
		for index := int32(0); index < section.Size; index++ {
			value := strconv.FormatInt(section.Min+int64(index), 10)
			less := int(section.Length - int32(len(value)))
			sections = append(sections, fmt.Sprintf("%s%s", strings.Repeat("0", less), value))
		}
		return sections, nil
	})
}

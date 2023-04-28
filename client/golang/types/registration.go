/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	METADATA = "metadata"
	PROXY    = "proxy"
	SERVER   = "server"
	COMPLEX  = "complex"
)

type Registration[T any] struct {
	InstanceId  string            `index:"0" json:"instance_id" xml:"instance_id" yaml:"instance_id"`
	Name        string            `index:"5" json:"name" xml:"name" yaml:"name"`
	Kind        string            `index:"10" json:"kind" xml:"kind" yaml:"kind"`
	Address     string            `index:"15" json:"address" xml:"address" yaml:"kind"`
	Content     T                 `index:"20" json:"content" xml:"content" yaml:"content"`
	Timestamp   int64             `index:"25" json:"timestamp" xml:"timestamp" yaml:"timestamp"`
	Attachments map[string]string `index:"30" json:"attachments" xml:"attachments" yaml:"attachments"`
}

type MetadataRegistration Registration[*Metadata]

func (that *MetadataRegistration) Any() *Registration[any] {
	return &Registration[any]{
		InstanceId:  that.InstanceId,
		Name:        that.Name,
		Kind:        that.Kind,
		Address:     that.Address,
		Content:     that.Content,
		Timestamp:   that.Timestamp,
		Attachments: that.Attachments,
	}
}

func (that *MetadataRegistration) InferService() []*Service {
	if nil == that.Content {
		return nil
	}
	return that.Content.Services
}

type MetadataRegistrations []*MetadataRegistration

func (that MetadataRegistrations) Of(kind string) MetadataRegistrations {
	var rs MetadataRegistrations
	for _, r := range that {
		if r.Kind == kind {
			rs = append(rs, r)
		}
	}
	return rs
}

func (that MetadataRegistrations) InferService() []*Service {
	var services []*Service
	for _, r := range that {
		if nil == r.Content {
			continue
		}
		for _, service := range r.Content.Services {
			services = append(services, service)
		}
	}
	return services
}

func (that MetadataRegistrations) InferReference() []*Reference {
	var references []*Reference
	for _, r := range that {
		if nil == r.Content {
			continue
		}
		for _, reference := range r.Content.References {
			references = append(references, reference)
		}
	}
	return references
}

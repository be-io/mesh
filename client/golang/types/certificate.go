/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Certificate struct {
	Csr string `index:"0" xml:"csr" json:"csr" yaml:"csr" comment:"Certification request"`
	Key string `index:"5" xml:"key" json:"key" yaml:"key" comment:"Certification private key"`
	Crt string `index:"10" xml:"crt" json:"crt" yaml:"crt" comment:"Certification"`
}

type Certificates struct {
	Root        *Certificate `index:"0" xml:"root" json:"root" yaml:"root" comment:"Root certification"`
	RootName    string       `index:"5" xml:"root_name" json:"root_name" yaml:"root_name" comment:"Root certification name"`
	Privacy     *Certificate `index:"10" xml:"privacy" json:"privacy" yaml:"privacy" comment:"Privacy certification"`
	PrivacyName string       `index:"15" xml:"privacy_name" json:"privacy_name" yaml:"privacy_name" comment:"Privacy certification"`
	Cert        *Certificate `index:"20" xml:"cert" json:"cert" yaml:"cert" comment:"Certification"`
	CertName    string       `index:"25" xml:"cert_name" json:"cert_name" yaml:"cert_name" comment:"Certification name"`
}

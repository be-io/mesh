/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Credential struct {
	ClientId  string `index:"0" json:"client_id" xml:"client_id" yaml:"client_id" comment:"OAuth client id"`
	ClientKey string `index:"5" json:"client_key" xml:"client_key" yaml:"client_key" comment:"OAuth client key"`
	Username  string `index:"10" json:"username" xml:"username" yaml:"username" comment:"OAuth authorize username"`
	Password  string `index:"15" json:"password" xml:"password" yaml:"password" comment:"OAuth authorize password"`
}

type AccessToken struct {
	Token        string `index:"0" json:"token" xml:"token" yaml:"token" comment:"OAuth access token"`
	Type         string `index:"5" json:"type" xml:"type" yaml:"type" comment:"OAuth access token type"`
	ExpiresAt    int64  `index:"10" json:"expires_at" xml:"expires_at" yaml:"expires_at" comment:"OAuth access token expires time in mills"`
	Scope        string `index:"15" json:"scope" xml:"scope" yaml:"scope" comment:"OAuth access token scope"`
	RefreshToken string `index:"20" json:"refresh_token" xml:"refresh_token" yaml:"refresh_token" comment:"OAuth access token refresh token"`
}

type AccessGrant struct {
	Code string `index:"0" json:"code" xml:"code" yaml:"code" comment:"OAuth grant code"`
}

type AccessCode struct {
	Code        string `index:"0" json:"code" xml:"code" yaml:"code" comment:"OAuth authorize code"`
	RedirectURI string `index:"5" json:"redirect_uri" xml:"redirect_uri" yaml:"redirect_uri" comment:"OAuth authorize redirect uri"`
}

type AccessID struct {
	ClientId  string `index:"0" json:"client_id" xml:"client_id" yaml:"client_id" comment:"OAuth client id"`
	UserId    string `index:"5" json:"user_id" xml:"user_id" yaml:"user_id" comment:"OAuth user id"`
	ExpiresAt int64  `index:"10" json:"expires_at" xml:"expires_at" yaml:"expires_at" comment:"OAuth access token expires time in mills"`
}

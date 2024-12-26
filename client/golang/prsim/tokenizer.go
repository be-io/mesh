/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
	"time"
)

var ITokenizer = (*Tokenizer)(nil)

// Tokenizer spi
//
// Trusts:
//
// OAuth2:
//
//	+--------+                               +---------------+
//	|        |--(A)- Authorization Request ->|   Resource    |
//	|        |                               |     Owner     |
//	|        |<-(B)-- Authorization Grant ---|               |
//	|        |                               +---------------+
//	|        |
//	|        |                               +---------------+
//	|        |--(C)-- Authorization Grant -->| Authorization |
//	| Client |                               |     Server    |
//	|        |<-(D)----- Access Token -------|               |
//	|        |                               +---------------+
//	|        |
//	|        |                               +---------------+
//	|        |--(E)----- Access Token ------>|    Resource   |
//	|        |                               |     Server    |
//	|        |<-(F)--- Protected Resource ---|               |
//	+--------+                               +---------------+
//
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Tokenizer interface {

	// Apply
	// @MPI("mesh.trust.apply")
	Apply(ctx context.Context, kind string, duration time.Duration) (string, error)

	// Verify
	// @MPI("mesh.trust.verify")
	Verify(ctx context.Context, token string) (bool, error)

	// Quickauth OAuth2 quick authorize, contains grant code and code authorize.
	// @MPI("mesh.oauth2.quickauth")
	Quickauth(ctx context.Context, credential *types.Credential) (*types.AccessToken, error)

	// Grant OAuth2 code grant.
	// @MPI("mesh.oauth2.grant")
	Grant(ctx context.Context, credential *types.Credential) (*types.AccessGrant, error)

	// Accept OAuth2 accept grant code.
	// @MPI("mesh.oauth2.accept")
	Accept(ctx context.Context, code string) (*types.AccessCode, error)

	// Reject OAuth2 reject grant code.
	// @MPI("mesh.oauth2.reject")
	Reject(ctx context.Context, code string) error

	// Authorize OAuth2 code authorize.
	// @MPI("mesh.oauth2.authorize")
	Authorize(ctx context.Context, code string) (*types.AccessToken, error)

	// Authenticate OAuth2 authenticate.
	// @MPI("mesh.oauth2.authenticate")
	Authenticate(ctx context.Context, token string) (*types.AccessID, error)

	// Refresh OAuth2 auth token refresh.
	// @MPI("mesh.oauth2.refresh")
	Refresh(ctx context.Context, token string) (*types.AccessToken, error)
}

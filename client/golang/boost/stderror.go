//go:build !unix && !windows

/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"github.com/opendatav/mesh/client/golang/log"
	"runtime"
)

func RedirectStderrFile(ctx context.Context, stderr string) {
	log.Error(ctx, "Redirect stderr not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
}

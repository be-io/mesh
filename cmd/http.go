/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	Provide(new(HTTP))
}

type HTTP struct {
}

func (that *HTTP) Home(ctx context.Context) *cobra.Command {
	var addr, node, method, data, location, output string
	var headers, parts []string
	var insecure bool
	h2 := &cobra.Command{
		Use:     "http",
		Aliases: []string{"h1"},
		Version: prsim.Version,
		Short:   "Mesh curl.",
		Long:    "Mesh curl.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if len(args) < 1 {
				log.Info(mtx, "Mesh h1 command must use valid arguments. ")
				return
			}
			uri, err := types.FormatURL(args[len(args)-1])
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			request, err := func() (*http.Request, error) {
				if "" != data {
					headers = append(headers, "Content-Type: application/json")
					return http.NewRequestWithContext(ctx, method, uri.String(), bytes.NewBufferString(data))
				}
				form := url.Values{}
				headers = append(headers, "Content-Type: application/x-www-form-urlencoded")
				return http.NewRequestWithContext(ctx, method, uri.String(), strings.NewReader(form.Encode()))
			}()
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			for _, header := range headers {
				kv := strings.SplitN(header, ":", 2)
				if len(kv) > 1 {
					request.Header.Add(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
				}
			}
			c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}}}
			resp, err := c.Do(request)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			defer func() {
				log.Catch(resp.Body.Close())
			}()
			o, err := io.ReadAll(resp.Body)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, string(o))
		},
	}
	h2.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	h2.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	h2.Flags().StringVarP(&method, "request", "X", http.MethodGet, "Specify request command to use")
	h2.Flags().StringVarP(&data, "data", "i", "", "HTTP POST data")
	h2.Flags().StringVarP(&location, "location", "L", "", "Follow redirects")
	h2.Flags().StringVarP(&output, "output", "o", "", "Download output location")
	h2.Flags().BoolVarP(&insecure, "insecure", "k", false, "Allow insecure server connections when using SSL")
	h2.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "Pass custom header(s) to server")
	h2.Flags().StringArrayVarP(&parts, "form", "F", []string{}, "Specify multipart MIME data")
	return h2
}

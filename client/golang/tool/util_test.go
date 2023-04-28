/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/types"
	"go/build"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestFido(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:7304/fedoServer/req/getToken", bytes.NewBufferString("{\"nonce\":\"C0A86ED11661593110588821504\",\"sign\":\"CB483797D021032FB30D22BCB6724F33\",\"req\":\"{}\",\"encryption\":\"NONE\",\"fromAppId\":\"566f721b-14d7-4da9-b7a7-f01717f0d52f\"}"))
	if nil != err {
		t.Error(err.Error())
		return
	}
	r.Header.Set("mesh-proto", "native")
	r.Header.Set("mesh-urn", "edges.net.mesh.0001000000000000000000000000000000000.jg0100004000000000.trustbe.cn")
	r.Header.Set("Content-Type", "application/json")
	s, err := Client.Do(r)
	if nil != err {
		t.Error(err.Error())
		return
	}
	b := &bytes.Buffer{}
	if _, err = io.Copy(b, s.Body); nil != err {
		t.Error(err.Error())
		return
	}
	t.Log(b.String())
}

func TestHash(t *testing.T) {
	t.Log(reflect.Value{} == reflect.Value{})
	t.Log(Hash("x.x.x"))
	t.Log(Hash("x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x"))
}

func TestAssert(t *testing.T) {
	comparator := func(k string, x map[string]any) bool {
		return nil != x && nil != x[k] && "" != x[k]
	}
	t.Log(comparator("1", nil))
	t.Log(comparator("1", map[string]any{"1": 1}))
	t.Log(comparator("1", map[string]any{"1": nil}))
	t.Log(comparator("1", map[string]any{"1": time.Time{}}))
	t.Log(comparator("1", map[string]any{}))
	t.Log(comparator("1", map[string]any{"1": ""}))
	t.Log(comparator("1", map[string]any{"1": map[string]any{}}))
}

func TestCopyArray(t *testing.T) {
	type X struct {
		X time.Duration
		Y types.Time
		Z types.Duration
	}
	x := &X{X: time.Second, Y: types.Time(time.Now()), Z: types.Duration(time.Second)}
	b, err := json.Marshal(x)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(b))

	var xx X
	if err = json.Unmarshal(b, &xx); nil != err {
		t.Error(err)
		return
	}
	t.Log(xx.X.Nanoseconds(), time.Time(xx.Y).UnixMilli(), time.Duration(xx.Z).Milliseconds())

	addrs := NewAddrs("1.1.1.1:80;1.1.1.1:81;1.1.1.1:82;1.1.1.1:83", 0)
	t.Log(addrs.Many())
}

func TestRegexReplace(t *testing.T) {
	for _, q := range []string{
		"a.b.c:80?p=1.1.1.1:90",
		"a.b.c?p=http://1.1.1.1:90/",
		"a.b.c:80?p=http://1.1.1.1:90/",
		"a.b.c:80",
		"a.b.c:80?p=http://1.1.1.1:90/",
		"https://a.b.c:80?p=http://1.1.1.1:90/",
		"h2c://a.b.c:80?p=http://1.1.1.1:90/",
		"192.168.10.49:7304?p=https://7.16.0.1:6099",
		"sockets5://192.168.10.49:7304?p=https://7.16.0.1:6099",
	} {
		u, err := types.FormatURL(q)
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(u.String())
	}
	rules := map[string]string{
		"href=\"/": "href=\"./",
		"src=\"/":  "src=\"./",
		"href='/":  "href='./",
		"src='/":   "src='./",
	}
	for r, p := range rules {
		if pattern, err := regexp.Compile(r); nil != err {
			t.Error(err)
		} else {
			t.Log(string(pattern.ReplaceAll([]byte("href=\"/a/b/c\" src=\"/a/b/c\" href='/a/b/c' src='/a/b/c'"), []byte(p))))
		}
	}

	var patterns = []string{
		"href=\"/",
		"src=\"/",
		".href = \"/",
		".src = \"/",
		"baseAssets:\"/\"",
		"apiEndpoint:\"/",
		"projectEndpoint:\"/",
		"proxyEndpoint:\"/",
		"wsEndpoint:\"/",
		"self\":\"http:\\/\\/[a-zA-Z0-9|\\.|-|_|:]+\\/",
		"self\":\"http:\\\\/\\\\/[a-zA-Z0-9|\\.|-|_|:]+\\\\/",
	}

	exp := &bytes.Buffer{}
	for index, pattern := range patterns {
		exp.WriteString(fmt.Sprintf("%s%s", Ternary(0 == index, "", "|"), pattern))
	}
	pattern, err := regexp.Compile(exp.String())
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(pattern.ReplaceAll([]byte(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta id="viewport" name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="/assets/images/logos/favicon.ico">
    <title>Loading&hellip;</title>
    <!-- 1.6.52 -->
    <meta name="description" content="">
    <base href="/" />
    <link id = "vendor" rel="stylesheet" href="/assets/vendor.css">
    <link id="theme" rel="stylesheet">
    
  </head>
  <body>
    
    <script src="/assets/vendor-877dcd270ab4f697ea387c3414a91a71.js"></script>
    <script src="/assets/ui-c4b9150cbb6a0909f129ebc9c4a06ae7.js"></script>
    <div id="ember-basic-dropdown-wormhole"></div>
  </body>
</html>
`), []byte("${0}x/y/z/"))))
	t.Log(string(pattern.ReplaceAll([]byte("links\":{\"self\":\"http:\\/\\/127.0.0.1:7304\\/v2-beta\\/setting\"}"), []byte("${0}x/y/z/"))))
	t.Log(string(pattern.ReplaceAll([]byte("href=\"/\nsrc=\"/\n.href = \"/\n.src = \"/\nbaseAssets:\"/\"\napiEndpoint:\"/\nprojectEndpoint:\"/\nproxyEndpoint:\"/\nwsEndpoint:\"/\n"), []byte("${0}x/y/z/"))))
}

func TestAddressable(t *testing.T) {
	ctx := macro.Context()
	t.Log(Addressable(ctx, "127.0.0.1"))
	t.Log(Addressable(ctx, "127.0.0.1:90"))
	t.Log(Addressable(ctx, "https://127.0.0.1"))
	t.Log(Addressable(ctx, "a.b.c"))
	t.Log(Addressable(ctx, "a.b.c:80"))
	t.Log(Addressable(ctx, "https://a.b.c:80"))
	t.Log(Addressable(ctx, "https://a.b.c:80?v=1.5.6"))

	t.Log(Addressable(ctx, "10.99.107.33:7304"))
	t.Log(Addressable(ctx, "https://10.99.107.33:7304"))
	t.Log(Addressable(ctx, "h2c://10.99.107.33:7304"))
	t.Log(Addressable(ctx, "h2c://10.99.107.33:7304?v=1.5.6"))
}

func TestVersionCompare(t *testing.T) {
	v := &types.Versions{Version: "1.5.5.3"}
	t.Log(v.Compare("1.5.*"))
}

func TestURLQuery(t *testing.T) {
	for _, uname := range []string{"https://abc.com:9900?v=1.5.6", "https://10.99.1.33:80?v=1.5.6", "10.99.1.33?v=1.5.6.1", "10.99.1.33?v=", "10.99.1.33?", "10.99.1.33"} {
		uri, err := url.Parse(uname)
		if nil != err {
			t.Error(err)
			continue
		}
		version := &types.Versions{Version: uri.Query().Get("v")}
		t.Log(uri.Hostname(), version.Compare("1.5.6"))
	}
	for _, addr := range []string{"https://abc.com:9900?v=1.5.6", "abc.com:9900?v=1.5.6", "abc.com:9900"} {
		uri, err := url.Parse(Ternary(strings.Contains(addr, "://"), addr, fmt.Sprintf("https://%s", addr)))
		if nil != err {
			t.Error(err)
			continue
		}
		t.Log(uri.Query().Get("v"))
	}
}

func TestCompare(t *testing.T) {
	zeroCompare := func(x interface{}) bool {
		if nil == x || "" == x {
			return true
		}
		return false
	}
	t.Log(zeroCompare(nil))
	t.Log(zeroCompare(""))
	t.Log(zeroCompare(time.Time{}))
}

func TestNewTraceID(t *testing.T) {
	t.Log(IPHex.Get())
	t.Log(NextID())
	t.Log(NewTraceId())
}

func TestTime(t *testing.T) {
	var ti types.Time
	if err := codec.Jsonizer.Unmarshal([]byte("1658317077473"), &ti); nil != err {
		t.Error(err)
		return
	}
	x, err := codec.Jsonizer.Marshal(ti)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(x))
}

func TestAnyone(t *testing.T) {
	var zero map[string]string
	nz := map[string]string{"1": "1"}
	t.Log(Anyone("", "", "1"))
	t.Log(Anyone(zero, nz, nz))
}

func TestCodec(t *testing.T) {
	var x map[string]string
	if err := codec.Jsonizer.Unmarshal([]byte("{\"x\":\"y\"}"), &x); nil != err {
		t.Error(err)
		return
	}
	t.Log(x)
	y := &map[string]string{}
	if err := codec.Jsonizer.Unmarshal([]byte("{\"x\":\"y\"}"), y); nil != err {
		t.Error(err)
		return
	}
	t.Log(*y)
}

func TestReplace(t *testing.T) {
	pattern := regexp.MustCompile(`tcp\(([\w|\W]+)\)`)
	str := pattern.ReplaceAllString("omega:Gaia@123@tcp(ido-mysql-headless:3306)/metabase", "$1")
	t.Log(str)
}

func TestRegexSub(t *testing.T) {
	reg, err := regexp.Compile(fmt.Sprintf("([-a-zA-Z\\d]+)\\.%s\\.([-a-zA-Z\\d]+)\\.((?i)JG[-a-zA-Z\\d]{2}%s[-a-zA-Z\\d]{8}|(?i)LX[-a-zA-Z\\d]{6}%s[-a-zA-Z\\d]{1}%s)\\.%s\\.(net|cn|com)", "mesh", "000001", "000001", "|(?i)JG0100000x00000000|(?i)LX00000100000z0", "trustbe"))
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(reg.String())
	t.Log(reg.MatchString("xxx.mesh.xxx.jg0100000100000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.JG0100000100000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.JG0100000110000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.jg0100000x00000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.lx00000100000z0.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.JG0100000x00000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.LX00000100000z0.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.Jg0100000x00000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.Lx00000100000z0.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.lx0000010000010.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.LX0000010000010.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.lx0000020000010.trustbe.com"))
	t.Log(reg.MatchString("xxx.xxx.xxx.trustbe.com"))
	t.Log()
	t.Log(reg.MatchString("xxx.mesh.xxx.lx0000010000020.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.lx00000100000200.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.jg0100000200000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.jg0100001000000000.trustbe.com"))
	t.Log(reg.MatchString("xxx.mesh.xxx.jg0100000010000000.trustbe.com"))
	t.Log()
	aeg, err := regexp.Compile("([-a-zA-Z\\d]+)\\.([-a-zA-Z\\d]{4}000001[-a-zA-Z\\d]{8})\\.trustbe\\.(net|cn|com)")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(aeg.String())
	t.Log(aeg.MatchString("xxx.xxx.xxx.trustbe.com"))
	t.Log(aeg.MatchString("xxx.mesh.xxx.jg0100000100000000.trustbe.com"))
	t.Log(aeg.MatchString("xxx.mesh.xxx.JG0100000100000000.trustbe.com"))
	t.Log(aeg.MatchString("xxx.mesh.xxx.jg0100000200000000.trustbe.com"))
	t.Log(aeg.MatchString("xxx.mesh.xxx.jg0100001000000000.trustbe.com"))
	t.Log(aeg.MatchString("xxx.mesh.xxx.jg0100000010000000.trustbe.com"))
	leg, err := regexp.Compile(fmt.Sprintf("([-a-zA-Z\\d]+)\\.%s\\.([-a-zA-Z\\d]+)\\.((?i)JG[-a-zA-Z\\d]{2}%s[-a-zA-Z\\d]{8}|(?i)LX[-a-zA-Z\\d]{6}%s[-a-zA-Z\\d]{1}%s)\\.%s\\.(net|cn|com)", "mesh", "000001", "000001", "|(?i)JG0100000x00000000|(?i)LX00000100000z0", "trustbe"))
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(leg.String())
	t.Log(leg.MatchString("weave.net.mesh.0001000000000000000000000000000000000.lx0000010000010.trustbe.cn"))
}

func TestRegex(t *testing.T) {
	reg, err := regexp.Compile("([-a-z0-9]+)\\.zzz\\.([-a-z0-9]+).(LX00000001|lx00000001)\\.trustbe\\.(net|cn)")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(reg.MatchString("xxx.xxx.xxx.trustbe.com"))
	t.Log(reg.MatchString("xxx.xxx.xxx.trustbe.cn"))
	t.Log(reg.MatchString("xxx.xxx.000.trustbe.cn"))
	t.Log(reg.MatchString("xxx.xxx.000.trustbe.cn"))
	t.Log(reg.MatchString("xxx.zzz.000.Lx00000001.trustbe.cn"))
	t.Log(reg.MatchString("xxx.xxx.000.LX00000001.trustbe.cn"))
	t.Log(reg.MatchString("xxx.zzz.000.LX00000001.trustbe.cn"))
	t.Log(reg.MatchString("zz.bb.11.xxx.xxx.000.lx00000001.trustbe.cn"))
	t.Log(reg.MatchString("xxx.xxx.000.lx00000001.trustbe.cn"))

	e, err := regexp.Compile("([-a-z0-9]+)\\.(asset|theta|edge)\\.([-a-z0-9]+)\\.(LX0000010000010|JG0100000100000000|lx0000010000010|jg0100000100000000)\\.trustbe\\.(net|cn|com)")
	if nil != err {
		t.Log(err)
		return
	}
	t.Log(e.MatchString("explore.data.api.asset.0000000000000000000000000000000000000.jg0100000100000000.trustbe.cn"))

	x, err := regexp.Compile("([-a-z0-9]+)\\.mesh\\.([-a-z0-9]+)\\.(LX0000010000060|JG0100000600000000|lx0000010000060|jg0100000600000000|LX0000000000000|JG0000000000000000|lx0000000000000|jg0000000000000000)\\.trustbe\\.(net|cn)")
	if nil != err {
		t.Log(err)
		return
	}
	t.Log(x.MatchString("query_list.task_offline.prom.mesh.0001000000000000000000000000000000000.lx0000000000000.trustbe.cn"))

	z, err := regexp.Compile("([-a-z0-9]+)\\.(LX0000010000010|JG0100000100000000|lx0000010000010|jg0100000100000000)\\.trustbe\\.(net|cn)")
	if nil != err {
		t.Error(err)
		return
	}
	for index := 0; index < 10; index++ {
		t.Log(z.MatchString("all.find.cooperation.omega.0000000000000000000000000000000000080.lx0000010000010.trustbe.cn"))
	}
	v, err := regexp.Compile("A(.*)")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(v.MatchString("A.xxx"))
	t.Log(v.MatchString("A.x.x.x"))
	t.Log(v.MatchString("A.x.x.X"))
	t.Log(v.MatchString("a.x.x.X"))
	t.Log(v.MatchString("ax.x.X"))
	t.Log(v.MatchString("A"))
	t.Log(v.MatchString("a"))

	n, err := regexp.Compile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(n.MatchString("127.0.0.1"))
	t.Log(n.MatchString("127.0.0"))
	t.Log(n.MatchString("255.255.255.255"))
	t.Log(n.MatchString("255.255.255.255."))
	t.Log(n.MatchString("255.255.255.255:80"))

	m, err := regexp.Compile("// @([a-zA-Z0-9]+(\\(.*\\))?)")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(m.MatchString("@MPS(1)"))
	t.Log(m.MatchString("// @MPS"))
	t.Log(m.MatchString("// @MPS(0)"))
	t.Log(m.MatchString("// @MPS(x=\"1\")"))
	t.Log(m.MatchString("// @MPS()"))
	t.Log(m.FindAllString(`
// FullType returns the fully qualified type of e.
// Examples, assuming package net/http:
// FullType(int) => "int"
// FullType(Handler) => "http.Handler"
// FullType(io.Reader) => "io.Reader"
// FullType(*Request) => "*http.Request"
// @MPI(name="x",x="2", z=3,    n="8", l=true)   the fully qualified type of e. 12222313
// @MPI(name="x",x="2", z=4,    n="8", l=true)   the fully qualified type of e. 12222313
	`, 3))
	t.Log(m.FindAllString(`
// @MPI("mesh.graph.mutate")//Backup(ctx context.Context, in *types.BackupRequest) (*types.BackupResponse, error)
//Restore(ctx context.Context, in *RestoreRequest, opts ...grpc.CallOption) (*Status, error)
//Export(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
//ReceivePredicate(ctx context.Context, opts ...grpc.CallOption) (Worker_ReceivePredicateClient, error)
//MovePredicate(ctx context.Context, in *MovePredicatePayload, opts ...grpc.CallOption) (*api.Payload, error)
//Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Worker_SubscribeClient, error)
//UpdateGraphQLSchema(ctx context.Context, in *UpdateGraphQLSchemaRequest, opts ...grpc.CallOption) (*UpdateGraphQLSchemaResponse, error)
//DeleteNamespace(ctx context.Context, in *DeleteNsRequest, opts ...grpc.CallOption) (*Status, error)
//TaskStatus(ctx context.Context, in *TaskStatusRequest, opts ...grpc.CallOption) (*TaskStatusResponse, error) 
	`, 31))

}

func TestImport(t *testing.T) {
	pkg, err := build.Import("github.com/be-io/mesh/client/golang/prsim", "", build.ImportComment)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(pkg)
}

func TestURIParse(t *testing.T) {
	uri, err := url.Parse("mysql://x:y@127.0.0.1:3306/z?v=k&m=n")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(uri)
	t.Logf("%v@tcp(%s)%s?%s", uri.User, uri.Host, uri.Path, uri.RawQuery)

	uri, err = url.Parse("x:y@tcp(127.0.0.1:3306)/z?v=k&m=n")
	if nil != err {
		t.Error(err)
	}
	t.Logf("%v@tcp(%s)%s?%s", uri.User, uri.Host, uri.Path, uri.RawQuery)
}

func TestFormatLicense(t *testing.T) {
	buff, err := codec.Jsonizer.Marshal(map[string]string{
		"license": `-----BEGIN TRUSTBE LICENSE V0-----
H4sIAAAAAAAA/wAACv/1uA66ZN24/EERn+Hd2AHFvZ5hn9QCbC/S8mvpXxlhSXnY
VomxdNUTsItBbL1/hz48Pr/Aqk4Dh0YbOD68XxcBo333lx5rC7/fskLyreQrWCeB
lLEmXmkzIRRCQ1g8EcPvbDlogOIrQ9FmDX47OZx7Ttml3haUlw8vE5Zi1SKw4PtT
I45KV7YUgwDtFxpgCEIH+Af5Za2/0ozscRQQtGQ7iL1lr2hoftPXtE2hNYypV4iD
RwsvmnBg2Vg6aDHDf6arEtrcaPMUn/eljwsPnsJuz556hBhjsatPPdk977j662pu
zlGBjT9M10pqB1eB4VlulmHfqMYRSDrfPKrSMOcnYF++H0LjQg0tmWK8f6DhaqA2
Gzy5E3/9ay8QeWzTl1hz/Nnu59x2Bzq3fS9HfnFR0pvCbWjTt9vM+e2qpMSIbW0K
oT2msWpwWgrANsrPCgQayJC/PkQOQyMM6qy1qcYZwkJ7nYeDu++QXHOCQFUlsA6U
EwFlYDnViR6oD+V/wvvlul3L6j5gbGZ1895BHD7wdF+CXB95p+/SfRFbms235CJY
DvPj8NBFZzHo8Rd4AaiRgQD5aXLScGdQb3vGUdU17Ehv7X+0XhyyY1wQwYRANdSy
1Ehq7xWZc7ymTX+zh4HD06f+tjpnLhxYO/MKBHg2WlFSxo6Zl1lmawhoEjSaRUUe
Dio4ldUh9K0xfEtWe7SDkQPWyuGSc9/vhrlODxEYw7Y0kNUUZKzB/KV3X12FCobx
3tOYBxx2VYeT1+m/W3/F7f7xR/qm3iTOuL28kkP/B2Lcg5YHHviQTnfTmI0P5YWr
bGf8/9Jks867hC8V3nXSc+67Hlc4s4qi81S/LCT+tguVTKdzFzBEKkMtEVW92T93
J0rKwZw05EYlCi48SIoM9jjCyJ2XKE3Z4D5tFBWNKjYon823vxktg3IEA4thjg20
durH2y5/ZbPDQpRkYJMF0EpwUMjruPqHrynObmMhrHT2LaDzVV4TcwH8ijixGTaa
8LGu48HOvME185QsfHxRUQZFMFDQzmV9sKafGtw9lUaLOaN8A51NrrhDGW/si4+7
O8XUvuo7YtRYC9Z2wYd+jf90GgjBHbXHuSREh+j1+XXASJe+UxZtAwW+f3akJ3J1
NXQZ3GC3V5eBsKdqAeAOqSnr+KTNj6JLOVJ1UWudjlZnLY+COOG0JlIuoOTR6x13
NIKtIWhie3GDTdD3/jaR1UDyQ3p4PPDgQ2IOmQpgZZHcyuacG1+Q9fjbYzGXPmgT
GjmVy84KeGZL7CV0YSWhtF2wnjWsrYCaN6yEkU1afb/+AKTQr9M/wTXMfco2rqhG
d1INFVQpdK03aa2bOcbx529PZhQKeBg61Dzs1Cx12jWTP04uiz4/RgPuyCocEh0p
W3hvRBxJQk6Aq99lhjuRWj0K37ogZp/DwOVv+KTr/AVddc38q+F1K9ThX4nksAo6
StlhbaJg5YNqYY1CiDBeCm0vupsMfNT2hEYU+rxObFZA6t4hOOgsTc8wghAqRAUY
J3dTkrON+r/NQREL0/NP3xg3O0YYN1UyQT0wrAv3iEq9EfpMjwcbyreSKYUrEVHH
l4SBtusWPBbmAQq2HADLq04mUJk+WGDY1FnTpLyqbTjQpb414BiBmcpEM1GAosTX
07dl/vn3shQr24rpHp90T2vCxUqE5sQqsN99++8gvrvMEwVIedeHS3Mp8pcVUBOa
mh6To/DE8frQQvDqRj35KUiBgDF3Y4inA8a7cZE6ExwpBGo96QOptcGbW3IJ7TnB
/yNkRSPx21KCR+uhMfFJphjajRVlWnTm/WS+cykVHWTlQdEI5NTcFWpfPAfFDpZk
L8wgnVpNW2cnJme5U1jbv9pB18GA/jJDzGGLPu46PaHE3dpPn7BxSyTHCXPA+VbM
4MvnpwAL5+SXOcrPMGTtqKZnFwrZltOTl7nccxwlrAH83f9RYN46Ie2Mr1nF2FM+
xF+ClqN34IQKX3xYJ8RpIdVlhNeltUCLizpBOVILheNJdMBCxau3/Ozbnq0wwyQy
n1DL7pOhMxuaKJMTkwbch69aZTT/n299C/M9BGq2so5mRZKxA/pvpEzIJs1G2iv6
7bFuRfWOt7LU0J/rpHkxYNMy1qBxX5/V1SwNetVFkfFeT66VxZTSu/cgsycAcEa1
QRWrDXVrgn/GwUjs75bgGLCVHJYIL+BPt83c+RfpOFLd5aYVeaFBJwgM+H0LSR5S
KPHOpDdX1oWDKsFmoJRDY9ZrEDhJxMU+yGiToI/8Mk+Ske8gaI5oFD1wKNbRqWSI
swtrbz7psPuaTp/C+kYHdR0QQ6fw4nsD5iDEaU5tjxXIg/PiJGx80JI5SxDEe5H9
p7Ol+Nl9S5RCxQ5Cafr11Q8dvRAOULu+btN7CybZsJPmNLSgFVWs13ogld/ESHPz
7ejhvMwabYz2iNAS61ERjXEJLoO+KcPfC7zBpWOgzWV8pcB4Eb80tAnTHifIrPRz
BAXgFcZeINtSi+1Wj4Zw0UUaxnU5BZ8d6Mt9e4DHckjUObbvMi0PbGU893jl3yWi
OszLJaco3UVjncognziDpImAg6dJiZkrxzJkbF7kdVZSJifMZV53AYaGfbuRGZvV
E5iaWpafM1lvr4LgPyGKGUTz3hK0FoXx0BIPTRD5OikygTsxrZMiegoaZXl83e6u
wouMYPxzERc+jBdcC2WgsgeKoUS1VzndAhFuRYsjZKuI38je2yILx689vhA3V49L
qxK5ozzkbQECyrRiRIqYnRtan/MDxYjHe279ZXXsFev+npechnrfxtcmVj6V22y8
OUTzt/vpXoxsmZ2281hKh8ZxFUk3DqxcHOTITKj0fx3CEK+sfOy68HJxfdd6beHx
XbI8CYZ8fJF4kCXJZtmvPKGBcw9/llV7liRMHmJdg6yZYBPhkm5OM0w4RD8THENR
H6h4E0rpculLqmPleVfeGOLz4+VO5DRBroMUDTsy/Dg19kifBBWAAbW8RmSDaEay
/AaUE0Tm/cBJipvnFXoPzV8P5jkQyY1n93Qr+AaMdkk+smozWtP+FakvyPbBzGjd
BlVSM36mvoijB0wN55JtNWag1kBhP9jzyQMtV+60lYc03Bq3PkQqGpZBYr6ogSih
VMzKPqPVaBc62zOGpfWTEOUVft5y0g7JE7zyBWqsxKUR8tbCslhf9BRqFMtxjz0t
8dU6cpbkJ+0fjyzbKJVBg8k5DWUnNrXZfMUhD01uJ+HW6V0htaJwmj+kkTLVuz5t
7vvPpQD3vSa/WWwM3aCDnYQBBS8ZVbwQ+Etv5dGR66u2HYJWkPyahrrhSArr+e+U
i1iJ+JbgqXgk3jsmaC3oockDM8HD20oWblTb8jxs/D7Z36niD+3OOPezzKkRF/e6
oLp2b1jI9e2Pn5Rpd/F7YMv9v6ExC+1tSsVMfwfhDQEAAP//2PbVmgAKAAA=
-----END TRUSTBE LICENSE V0-----`,
	})
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(buff))
}

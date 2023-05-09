package proxy

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequestURI(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://approvalSave.approval.serving.edge.0001000000000000000000000000000000000.jg0100000100000000.ducesoft.net/mesh-rpc/v1", nil)
	fmt.Printf("uri:%s\n", req.RequestURI)

	fmt.Printf("url:%s\n", req.URL.String())

	fmt.Printf("host:%s\n", req.Host)

	req.URL.Path = "approvalSave.approval.serving.pandora.0001000000000000000000000000000000000.jg0100000100000000.ducesoft.net/mesh-rpc/v1"
	fmt.Printf("host:%s\n", req.Host)

}

func TestPort(t *testing.T) {
	t.Log(7605 / 100 * 100)
}

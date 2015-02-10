package test

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFreePort(t *testing.T) {
	_, err := FreePort()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpServer(t *testing.T) {
	server, err := NewHttpServer(0)
	if err != nil {
		t.Fatal(err)
	}

	reqCh := server.AddResponse(http.StatusOK, http.Header{"Test-Header": {"test header"}}, []byte("test body"))

	req, err := http.NewRequest("GET", "http://"+server.Address()+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	rcvReq := <-reqCh

	if rcvReq.Method != req.Method {
		t.Error(rcvReq, req)
	} else if rcvReq.Host != req.Host {
		t.Error(rcvReq, req)
	} else if resp.StatusCode != http.StatusOK {
		t.Error(rcvReq, req)
	} else if resp.Header.Get("Test-Header") != "test header" {
		t.Error(rcvReq, req)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else if string(body) != "test body" {
		t.Error(string(body), "test body")
	}
}

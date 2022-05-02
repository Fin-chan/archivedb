package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type httpCli struct {
	http.Client

	header http.Header
}

func NewWithHeader(header map[string]string) (cli *httpCli, err error) {
	h := http.Header{}
	for k, v := range header {
		h.Set(k, v)
	}

	cli = &httpCli{
		Client: http.Client{},
		header: h,
	}
	return
}

func (h *httpCli) Post(url string, data string) ([]byte, error) {
	return h.do("POST", url, []byte(data))
}

func (h *httpCli) do(method, url string, data []byte) ([]byte, error) {
	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = h.header

	resp, err := h.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (h *httpCli) Get(url string) ([]byte, error) {
	return h.do("GET", url, nil)
}

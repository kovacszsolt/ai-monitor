package proxy

import (
	"net"
	"net/http"
	"time"
)

func NewHTTPClient(connectionTimeoutMS, requestTimeoutMS int) *http.Client {
	dialer := &net.Dialer{Timeout: time.Duration(connectionTimeoutMS) * time.Millisecond}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: time.Duration(requestTimeoutMS) * time.Millisecond,
		// Do not negotiate gzip with upstream; we decode explicitly when providers still send gzip.
		DisableCompression: true,
	}
	return &http.Client{Transport: transport}
}

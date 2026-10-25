package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func maybeDecodeUpstreamBody(hdr http.Header, body []byte) ([]byte, error) {
	if len(body) == 0 {
		return body, nil
	}
	if !shouldGunzipResponse(hdr, body) {
		return body, nil
	}
	decoded, err := gunzip(body)
	if err != nil {
		return body, fmt.Errorf("gunzip upstream body: %w", err)
	}
	clearEncodedContentHeaders(hdr)
	return decoded, nil
}

func shouldGunzipResponse(hdr http.Header, body []byte) bool {
	enc := strings.ToLower(strings.TrimSpace(hdr.Get("Content-Encoding")))
	if strings.Contains(enc, "gzip") {
		return true
	}
	return len(body) >= 2 && body[0] == 0x1f && body[1] == 0x8b
}

func gunzip(body []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}

func clearEncodedContentHeaders(hdr http.Header) {
	hdr.Del("Content-Encoding")
	hdr.Del("Content-Length")
}

package proxy

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"testing"
)

func TestMaybeDecodeUpstreamBodyGzipHeader(t *testing.T) {
	plain := []byte(`{"usage":{"prompt_tokens":1}}`)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write(plain)
	_ = zw.Close()

	hdr := http.Header{"Content-Encoding": []string{"gzip"}}
	out, err := maybeDecodeUpstreamBody(hdr, buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("decoded: %s", out)
	}
	if hdr.Get("Content-Encoding") != "" {
		t.Fatal("expected content-encoding removed")
	}
}

func TestMaybeDecodeUpstreamBodyMagicWithoutHeader(t *testing.T) {
	plain := []byte(`{"ok":true}`)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write(plain)
	_ = zw.Close()

	hdr := http.Header{}
	out, err := maybeDecodeUpstreamBody(hdr, buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatal("magic gzip without header")
	}
}

func TestMaybeDecodeUpstreamBodyPlainJSON(t *testing.T) {
	plain := []byte(`{"ok":true}`)
	hdr := http.Header{}
	out, err := maybeDecodeUpstreamBody(hdr, plain)
	if err != nil || !bytes.Equal(out, plain) {
		t.Fatal("plain json should pass through")
	}
}

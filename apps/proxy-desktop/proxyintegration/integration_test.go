package tests

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/app"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/export"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/proxy"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

func TestE2EChatCompletion(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var v map[string]interface{}
		json.Unmarshal(body, &v)
		if v["model"] != "Qwen3.8-27B-8bit" {
			t.Fatalf("expected upstream model rewrite, got %v", v["model"])
		}
		time.Sleep(25 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"x","object":"chat.completion","choices":[{"message":{"content":"hi"}}],"usage":{"prompt_tokens":3,"completion_tokens":5,"prompt_tokens_details":{"cached_tokens":0}}}`))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 0
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "dolgozatmester", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local-omlx", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "qwen-local", prov.ID, "Qwen3.8-27B-8bit", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys:      dbApp.Keys,
		Providers: dbApp.Providers,
		Routes:    dbApp.Routes,
		Writer:    dbApp.Writer,
		MaxBody:   cfg.Server.MaxRequestBytes,
		Capture:   cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	reqBody := `{"model":"qwen-local","messages":[{"role":"user","content":"Hello!"}],"stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	time.Sleep(400 * time.Millisecond)

	var buf bytes.Buffer
	if err := export.ExportJSONL(context.Background(), dbApp.DB, &buf, export.Filter{}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("qwen-local")) {
		t.Fatalf("export missing metadata: %s", buf.String())
	}
	if bytes.Contains(buf.Bytes(), []byte("Hello!")) {
		t.Fatal("prompt leaked to export/db path")
	}

	var inTok, outTok, uncached sql.NullInt64
	var eff, gen sql.NullFloat64
	var genEst int
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT input_tokens, output_tokens, uncached_input_tokens, effective_tps, generation_tps, generation_speed_estimated
		FROM ai_requests
		WHERE requested_model = 'qwen-local'
		ORDER BY started_at DESC LIMIT 1`).Scan(&inTok, &outTok, &uncached, &eff, &gen, &genEst)
	if err != nil {
		t.Fatal(err)
	}
	if !inTok.Valid || inTok.Int64 != 3 || !outTok.Valid || outTok.Int64 != 5 {
		t.Fatalf("unexpected token counts: in=%v out=%v", inTok, outTok)
	}
	if !uncached.Valid || uncached.Int64 != 3 {
		t.Fatalf("expected uncached_input_tokens=3, got %v", uncached)
	}
	if !eff.Valid || eff.Float64 <= 0 {
		t.Fatalf("expected positive effective_tps, got %v", eff)
	}
	if !gen.Valid || gen.Float64 <= 0 {
		t.Fatalf("expected estimated generation_tps for non-stream, got %v", gen)
	}
	if genEst != 1 {
		t.Fatalf("expected generation_speed_estimated=1, got %d", genEst)
	}
}

func TestE2ERequestTag(t *testing.T) {
	var upstreamTag string
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamTag = r.Header.Get("X-AI-Proxy-Tag")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "tag-proj", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "tag-model", prov.ID, "upstream", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	post := func(tagHeader string) {
		reqBody := `{"model":"tag-model","messages":[{"role":"user","content":"t"}],"stream":false}`
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		if tagHeader != "" {
			req.Header.Set("X-AI-Proxy-Tag", tagHeader)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	post("integration-tag")
	time.Sleep(400 * time.Millisecond)
	if upstreamTag != "" {
		t.Fatalf("tag must not be forwarded upstream, got %q", upstreamTag)
	}
	var tag string
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT tag FROM ai_requests WHERE requested_model = 'tag-model' ORDER BY started_at DESC LIMIT 1`).Scan(&tag)
	if err != nil {
		t.Fatal(err)
	}
	if tag != "integration-tag" {
		t.Fatalf("expected tag integration-tag, got %q", tag)
	}

	var buf bytes.Buffer
	if err := export.ExportJSONL(context.Background(), dbApp.DB, &buf, export.Filter{}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"tag":"integration-tag"`)) {
		t.Fatalf("export missing tag: %s", buf.String())
	}

	post("")
	time.Sleep(400 * time.Millisecond)
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT tag FROM ai_requests WHERE requested_model = 'tag-model' ORDER BY started_at DESC LIMIT 1`).Scan(&tag)
	if err != nil {
		t.Fatal(err)
	}
	if tag != "" {
		t.Fatalf("expected empty tag without header, got %q", tag)
	}
}

func TestE2EInFlightProgress(t *testing.T) {
	release := make(chan struct{})
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"late"}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Capture.StoreContent = true
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, _ := dbApp.Projects.Create(context.Background(), "prog", "")
	_, apiKey, _ := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	prov, _ := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	_, _ = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	done := make(chan struct{})
	go func() {
		reqBody := `{"model":"m1","messages":[{"role":"user","content":"wait"}],"stream":false}`
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
		}
		close(done)
	}()

	deadline := time.Now().Add(2 * time.Second)
	var sawStarted bool
	for time.Now().Before(deadline) {
		var progress string
		err := dbApp.DB.QueryRowContext(context.Background(), `
			SELECT progress FROM ai_requests WHERE requested_model = 'm1' ORDER BY started_at DESC LIMIT 1`).Scan(&progress)
		if err == nil && progress == "started" {
			sawStarted = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !sawStarted {
		t.Fatal("expected in-flight row with progress=started")
	}
	close(release)
	<-done
	time.Sleep(400 * time.Millisecond)

	var progress, state string
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT progress, request_state FROM ai_requests WHERE requested_model = 'm1' ORDER BY started_at DESC LIMIT 1`).Scan(&progress, &state)
	if err != nil {
		t.Fatal(err)
	}
	if progress != "finished" || state != "completed" {
		t.Fatalf("after response: progress=%s state=%s", progress, state)
	}
}

func TestE2EPromptCaptureError(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Capture.StoreContent = true
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, _ := dbApp.Projects.Create(context.Background(), "err-proj", "")
	_, apiKey, _ := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	prov, _ := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	_, _ = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	// Valid model but no messages — metadata captured, prompt extract fails.
	reqBody := `{"model":"m1","stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	time.Sleep(400 * time.Millisecond)

	var n int
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM ai_request_errors e
		INNER JOIN ai_requests r ON r.id = e.request_id
		WHERE r.requested_model = 'm1' AND e.error_code = 'messages_missing'`).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 messages_missing error row, got %d", n)
	}
}

func TestE2EGzipUpstreamBody(t *testing.T) {
	plain := []byte(`{"choices":[{"message":{"content":"hi"}}],"usage":{"prompt_tokens":2,"completion_tokens":1}}`)
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") != "" {
			t.Errorf("proxy should not forward Accept-Encoding, got %q", r.Header.Get("Accept-Encoding"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, _ = zw.Write(plain)
		_ = zw.Close()
		w.Write(buf.Bytes())
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "gzip-proj", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	reqBody := `{"model":"m1","messages":[{"role":"user","content":"x"}],"stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d body %s", resp.StatusCode, body)
	}
	if !bytes.Contains(body, []byte(`"content":"hi"`)) {
		t.Fatalf("client body not json: %s", body)
	}
	time.Sleep(400 * time.Millisecond)

	var usageStatus string
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT usage_status FROM ai_requests ORDER BY started_at DESC LIMIT 1`).Scan(&usageStatus)
	if err != nil {
		t.Fatal(err)
	}
	if usageStatus != "ok" {
		t.Fatalf("usage_status: %q", usageStatus)
	}
}

func TestE2EHTTPFileDump(t *testing.T) {
	const upstreamBody = `broken-json`
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(upstreamBody))
	}))
	defer mock.Close()

	dir := t.TempDir()
	dumpDir := filepath.Join(dir, "dumps")
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Capture.FileDump = config.FileDumpConfig{
		Enabled:      true,
		Directory:    dumpDir,
		MaxBodyBytes: 4096,
	}
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "dump-proj", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	reqBody := `{"model":"m1","messages":[{"role":"user","content":"dump"}],"stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	time.Sleep(500 * time.Millisecond)

	entries, err := os.ReadDir(dumpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 dump file, got %d", len(entries))
	}
	raw, err := os.ReadFile(filepath.Join(dumpDir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), upstreamBody) {
		t.Fatalf("dump missing upstream body: %s", raw)
	}
}

func TestE2ERawBodyOnParseError(t *testing.T) {
	const upstreamBody = `not-json-but-upstream-reply`
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(upstreamBody))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Capture.StoreUpstreamBodyOnError = true
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "raw-proj", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	reqBody := `{"model":"m1","messages":[{"role":"user","content":"hi"}],"stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	time.Sleep(400 * time.Millisecond)

	var usageStatus, rawBody string
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT r.usage_status, COALESCE(res.raw_body, '')
		FROM ai_requests r
		LEFT JOIN ai_request_responses res ON res.request_id = r.id
		ORDER BY r.started_at DESC LIMIT 1`).Scan(&usageStatus, &rawBody)
	if err != nil {
		t.Fatal(err)
	}
	if usageStatus != "parse_error" {
		t.Fatalf("usage_status: %q", usageStatus)
	}
	if rawBody != upstreamBody {
		t.Fatalf("raw_body: %q", rawBody)
	}
}

func TestE2EStoreContent(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(25 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"assistant reply"}}],"usage":{"prompt_tokens":2,"completion_tokens":3}}`))
	}))
	defer mock.Close()

	dir := t.TempDir()
	cfg := config.Default()
	cfg.Storage.DatabasePath = filepath.Join(dir, "proxy.db")
	cfg.Capture.StoreContent = true
	migrations := filepath.Join("..", "migrations")

	dbApp, err := app.Open(cfg, migrations)
	if err != nil {
		t.Fatal(err)
	}
	defer dbApp.DB.Close()

	p, err := dbApp.Projects.Create(context.Background(), "content-proj", "")
	if err != nil {
		t.Fatal(err)
	}
	_, apiKey, err := dbApp.Keys.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	prov, err := dbApp.Providers.Add(context.Background(), "local", providers.TypeOpenAICompatible, mock.URL, "", 5000, 30000)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbApp.Routes.Add(context.Background(), "m1", prov.ID, "upstream", routing.ProtocolOpenAIChat)
	if err != nil {
		t.Fatal(err)
	}

	h := proxy.NewHandler(proxy.Services{
		Keys: dbApp.Keys, Providers: dbApp.Providers, Routes: dbApp.Routes,
		Writer: dbApp.Writer, MaxBody: cfg.Server.MaxRequestBytes, Capture: cfg.Capture,
	})
	srv := httptest.NewServer(h.AuthMiddleware(serverTestHandler(h)))
	defer srv.Close()

	reqBody := `{"model":"m1","messages":[{"role":"user","content":"Question?"}],"stream":false}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	time.Sleep(400 * time.Millisecond)

	var lastUser, assistant string
	err = dbApp.DB.QueryRowContext(context.Background(), `
		SELECT p.last_user_text, r.assistant_text
		FROM ai_request_prompts p
		INNER JOIN ai_request_responses r ON r.request_id = p.request_id
		ORDER BY p.created_at DESC LIMIT 1`).Scan(&lastUser, &assistant)
	if err != nil {
		t.Fatal(err)
	}
	if lastUser != "Question?" {
		t.Fatalf("last_user_text: %q", lastUser)
	}
	if assistant != "assistant reply" {
		t.Fatalf("assistant_text: %q", assistant)
	}
}

func serverTestHandler(h *proxy.Handler) http.Handler {
	mux := http.NewServeMux()
	h.Register(&muxAdapter{m: mux})
	return mux
}

type muxAdapter struct{ m *http.ServeMux }

func (a *muxAdapter) Get(p string, fn http.HandlerFunc)  { a.m.HandleFunc(p, fn) }
func (a *muxAdapter) Post(p string, fn http.HandlerFunc) { a.m.HandleFunc(p, fn) }

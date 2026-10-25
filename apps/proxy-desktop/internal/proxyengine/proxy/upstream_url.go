package proxy

import "strings"

// joinUpstreamURL combines provider base_url with a proxy route path (e.g. /v1/chat/completions).
// If base_url already ends with /v1 and the path also starts with /v1, the duplicate segment is dropped
// so both http://host:8001 and http://host:8001/v1 work as documented.
func joinUpstreamURL(baseURL, endpoint string) string {
	base := strings.TrimRight(baseURL, "/")
	if endpoint == "" {
		return base
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	if strings.HasSuffix(base, "/v1") && (endpoint == "/v1" || strings.HasPrefix(endpoint, "/v1/")) {
		endpoint = endpoint[len("/v1"):]
		if endpoint == "" {
			endpoint = "/"
		}
	}
	return base + endpoint
}

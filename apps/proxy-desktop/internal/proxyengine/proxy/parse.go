package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
)

var errInvalidJSONBody = errors.New("invalid json body")

func extractModel(body []byte) (model string, stream bool, rewritten []byte, err error) {
	var loose map[string]interface{}
	if err := json.Unmarshal(body, &loose); err != nil {
		return "", false, nil, fmt.Errorf("%w", errInvalidJSONBody)
	}
	m, _ := loose["model"].(string)
	if m == "" {
		return "", false, nil, fmt.Errorf("model is required")
	}
	s := false
	if v, ok := loose["stream"].(bool); ok {
		s = v
	}
	return m, s, body, nil
}

func rewriteModel(body []byte, upstreamModel string) ([]byte, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	v["model"] = upstreamModel
	return json.Marshal(v)
}

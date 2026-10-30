package content

import "strings"

// PromptErrorCode maps extract failures to stable codes for ai_request_errors.
func PromptErrorCode(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "invalid json"):
		return "invalid_json"
	case strings.Contains(msg, "messages missing"):
		return "messages_missing"
	case strings.Contains(msg, "prompt missing"):
		return "prompt_missing"
	case strings.Contains(msg, "input missing"):
		return "input_missing"
	default:
		return "prompt_extract_failed"
	}
}

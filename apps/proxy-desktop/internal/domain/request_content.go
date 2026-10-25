package domain

type AIRequestCaptureError struct {
	Stage        string
	ErrorCode    string
	ErrorMessage string
}

type AIRequestContent struct {
	RequestID      string
	RequestState   string
	ProxyError     string
	CaptureEnabled bool
	HasContent     bool

	LastUserText    string
	MessagesJSON    string
	PromptTruncated bool

	AssistantText      string
	ResponseTruncated  bool
	ToolCallsJSON      string
	ToolCallsTruncated bool
	RawBody            string
	RawBodyTruncated   bool

	CaptureErrors []AIRequestCaptureError
}

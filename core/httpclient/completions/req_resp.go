package completions

type CompletionRequest struct {
	Model       string           `json:"model"`
	Messages    []MessageRequest `json:"messages"`
	Temperature float64          `json:"temperature"`
	MaxTokens   int              `json:"max_tokens"`
	Stream      bool             `json:"stream"`
}

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Created   int64     `json:"created"`
	Model     string    `json:"model"`
	Choices   []Choice  `json:"choices"`
	Usage     Usage     `json:"usage"`
	Stats     Stats     `json:"stats"`
	ModelInfo ModelInfo `json:"model_info"`
	Runtime   Runtime   `json:"runtime"`
}

type Choice struct {
	Index        int     `json:"index"`
	Logprobs     *string `json:"logprobs"` // nullable
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Stats struct {
	TokensPerSecond  float64 `json:"tokens_per_second"`
	TimeToFirstToken float64 `json:"time_to_first_token"`
	GenerationTime   float64 `json:"generation_time"`
	StopReason       string  `json:"stop_reason"`
}

type ModelInfo struct {
	Arch          string `json:"arch"`
	Quant         string `json:"quant"`
	Format        string `json:"format"`
	ContextLength int    `json:"context_length"`
}

type Runtime struct {
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	SupportedFormats []string `json:"supported_formats"`
}

type CompletionError struct {
	Errors ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

func (e *CompletionError) Error() string {
	return e.Errors.Message
}
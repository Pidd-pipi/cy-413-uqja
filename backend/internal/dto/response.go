package dto

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type JournalPromptsResponse struct {
	Band    string   `json:"band"`
	Prompts []string `json:"prompts"`
}

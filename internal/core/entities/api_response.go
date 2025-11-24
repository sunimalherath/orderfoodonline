package entities

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Data    any    `json:"data,omitempty"`
}

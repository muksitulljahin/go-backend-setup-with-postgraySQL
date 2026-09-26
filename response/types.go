package response

// PaginationMeta defines metadata for list responses (Node.js Express style)
type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"totalPage"`
}

// APIResponse defines standard JSON format for all API responses
type APIResponse struct {
	Success    bool            `json:"success"`
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Data       interface{}     `json:"data,omitempty"`
	Meta       *PaginationMeta `json:"meta,omitempty"`
	Error      interface{}     `json:"error,omitempty"`
}

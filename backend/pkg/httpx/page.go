package httpx

// Pagination
type Pagination struct {
	Cursor int64 `json:"cursor,omitempty"`
	Limit  int64 `json:"limit,omitempty"`
	Page   int64 `json:"page,omitempty"`
	Total  int64 `json:"total,omitempty"`
}

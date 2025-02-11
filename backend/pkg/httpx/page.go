package httpx

// Pagination
type Pagination struct {
	Cursor int64 `json:"cursor,omitempty"`
	Limit  int64 `json:"limit,omitempty"`
	Page   int64 `json:"page,omitempty"`
	Total  int64 `json:"total,omitempty"`
}

// GetOffset returns the offset for the pagination.
func (p *Pagination) GetOffset() int64 {
	if p.Page > 0 {
		return (p.Page - 1) * p.Limit
	}

	return 0
}

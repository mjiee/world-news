package httpx

// Pagination
type Pagination struct {
	Cursor int `json:"cursor,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Page   int `json:"page,omitempty"`
	Total  int `json:"total,omitempty"`
}

// GetLimit returns the limit for the pagination.
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		return 20
	}

	return p.Limit
}

// GetOffset returns the offset for the pagination.
func (p *Pagination) GetOffset() int {
	if p.Page > 0 {
		return (p.Page - 1) * p.Limit
	}

	return 0
}

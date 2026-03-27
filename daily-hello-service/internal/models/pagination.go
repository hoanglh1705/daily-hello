package models

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *PaginationQuery) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PaginationQuery) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 100 {
		return 20
	}
	return p.Limit
}

func (p *PaginationQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type PaginatedResponse struct {
	Items interface{}    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

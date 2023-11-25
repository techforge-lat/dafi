package dafi

import "fmt"

type Pagination struct {
	Limit  uint
	Offset uint
}

func NewPagination(limit uint, offset uint) *Pagination {
	return &Pagination{Limit: limit, Offset: offset}
}

func (p *Pagination) SQL() string {
	if p.Limit == 0 {
		return ""
	}

	return fmt.Sprintf(" LIMIT %d OFFSET %d", p.Limit, p.Offset)
}

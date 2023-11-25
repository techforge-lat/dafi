package dafi

import "fmt"

type Criteria struct {
	Filter     Filter
	Sort       Sort
	Pagination Pagination
	err        error
}

func New(filter Filter, sort Sort, pag Pagination) *Criteria {
	return &Criteria{Filter: filter, Sort: sort, Pagination: pag}
}

func (c *Criteria) SQL() (string, []any) {
	where, args := c.Filter.SQL()
	c.err = c.Filter.Err()

	sort := c.Sort.SQL()
	pag := c.Pagination.SQL()

	return fmt.Sprintf("%s %s %s", where, sort, pag), args
}

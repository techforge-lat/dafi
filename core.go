package dafi

import "fmt"

type Criteria struct {
	Filter     Filter
	Sort       Sort
	Pagination Pagination
}

func New(filter Filter, sort Sort, pag Pagination) Criteria {
	return Criteria{Filter: filter, Sort: sort, Pagination: pag}
}

func (c Criteria) SQL() (string, []any, error) {
	where, args, err := c.Filter.SQL()
	if err != nil {
		return "", nil, err
	}

	sort := c.Sort.SQL()
	pag := c.Pagination.SQL()

	return fmt.Sprintf(" %s %s %s", where, sort, pag), args, nil
}

func (c Criteria) SQLWithReplacedNames(names map[string]string) (string, []any, error) {
	err := c.Filter.ReplaceAbstractNames(names)
	if err != nil {
		return "", nil, err
	}

	where, args, err := c.Filter.SQL()
	if err != nil {
		return "", nil, err
	}

	if err := c.Sort.ReplaceAbstractNames(names); err != nil {
		return "", nil, err
	}
	sort := c.Sort.SQL()

	pag := c.Pagination.SQL()

	return fmt.Sprintf(" %s %s %s", where, sort, pag), args, nil
}

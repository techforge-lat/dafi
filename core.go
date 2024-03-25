package dafi

type Criteria struct {
	Filter     Filter
	Sort       Sort
	Pagination Pagination
}

func New(filter Filter, sort Sort, pag Pagination) Criteria {
	return Criteria{Filter: filter, Sort: sort, Pagination: pag}
}

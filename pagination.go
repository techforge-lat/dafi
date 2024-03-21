package dafi

type Pagination struct {
	Limit  uint
	Offset uint
}

func NewPagination(limit uint, offset uint) Pagination {
	return Pagination{Limit: limit, Offset: offset}
}

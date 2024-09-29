package dafi

type SortType string

const (
	Asc  SortType = "asc"
	Desc SortType = "desc"
	None SortType = "none"
)

type SortBy string

type Sort struct {
	Field SortBy
	Type  SortType
}

type Sorts []Sort

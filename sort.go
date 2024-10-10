package dafi

type SortType string

const (
	Asc  SortType = "asc"
	Desc SortType = "desc"
	None SortType = ""
)

type SortBy string

type Sort struct {
	Field SortBy
	Type  SortType
}

type Sorts []Sort

func (s Sorts) IsZero() bool {
	return len(s) == 0
}

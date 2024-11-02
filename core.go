package dafi

const (
	groupOpenIndex = iota + 1
	fieldIndex
	operatorIndex
	valueIndex
	groupCloseIndex
	chainingKeyIndex
)

type Criteria struct {
	SelectColumns []string
	Joins         []string
	Filters       Filters
	Sorts         Sorts
	Pagination    Pagination
}

func New() Criteria {
	return Criteria{}
}

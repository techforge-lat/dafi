package dafi

type (
	FilterField string
	FilterValue any
)

type FilterOperator string

const (
	Equal          FilterOperator = "eq"
	NotEqual       FilterOperator = "ne"
	Greater        FilterOperator = "gt"
	GreaterOrEqual FilterOperator = "gte"
	Less           FilterOperator = "lt"
	LessOrEqual    FilterOperator = "lte"
	Like           FilterOperator = "like"
	In             FilterOperator = "in"
	NotIn          FilterOperator = "nin"
	Contains       FilterOperator = "contains"
	NotContains    FilterOperator = "ncontains"
	Is             FilterOperator = "is"
	IsNot          FilterOperator = "isn"
)

type FilterChainingKey string

const (
	And FilterChainingKey = "AND"
	Or  FilterChainingKey = "OR"
)

type Filter struct {
	IsGroupOpen   bool
	GroupOpenQty  int
	Field         FilterField
	Operator      FilterOperator
	Value         FilterValue
	IsGroupClose  bool
	GroupCloseQty int
	ChainingKey   FilterChainingKey
}

type Filters []Filter

func (f Filters) IsZero() bool {
	return len(f) == 0
}

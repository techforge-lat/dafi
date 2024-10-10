package dafi

type (
	FilterField string
	FilterValue any
)

type FilterOperator string

const (
	Equal          FilterOperator = "="
	NotEqual       FilterOperator = "!="
	Greater        FilterOperator = ">"
	GreaterOrEqual FilterOperator = ">="
	Less           FilterOperator = "<"
	LessOrEqual    FilterOperator = "<="
	Contains       FilterOperator = "CONTAINS"
	NotContains    FilterOperator = "NOT_CONTAINS"
	Is             FilterOperator = "IS"
	IsNot          FilterOperator = "IS_NOT"
	In             FilterOperator = "IN"
	NotIn          FilterOperator = "NOT_IN"
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

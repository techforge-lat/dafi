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
	Like           FilterOperator = "LIKE"
	NotLike        FilterOperator = "NOT_LIKE"
	ILike          FilterOperator = "ILIKE"
	NotILike       FilterOperator = "NOT_ILIKE"
	Is             FilterOperator = "IS"
	IsNull         FilterOperator = "IS_NULL"
	IsNot          FilterOperator = "IS_NOT"
	IsNotNull      FilterOperator = "IS_NOT_NULL"
	In             FilterOperator = "IN"
	NotIn          FilterOperator = "NOT_IN"
)

type FilterChainingKey string

const (
	And FilterChainingKey = "AND"
	Or  FilterChainingKey = "OR"
)

type Filter struct {
	IsGroupOpen  bool
	Field        FilterField
	Op           FilterOperator
	Value        FilterValue
	ChainingKey  FilterChainingKey
	IsGroupClose bool
}

type Filters []Filter

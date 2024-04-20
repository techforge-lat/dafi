package dafi

import (
	"errors"

	"github.com/techforge-lat/sqlcraft"
)

var (
	ErrInvalidFilterFormat = errors.New("invalid filter format, format is: field operator [value] chainingKey")
	ErrInvalidOperator     = errors.New("invalid operator, must be a =, <, >, <=, >=, <>, !=, IS, ILIKE, LIKE")
)

const (
	And = "AND"
	Or  = "OR"
)

const (
	Equal              = "="
	LessThan           = "<"
	GreaterThan        = ">"
	LessThanOrEqual    = "<="
	GreaterThanOrEqual = ">="
	NotEqual           = "<>"
	Is                 = "IS"
	IsNot              = "IS_NOT"
	ILike              = "ILIKE"
	NotILike           = "NOT_ILIKE"
	Like               = "LIKE"
	NotLike            = "NOT_LIKE"
	In                 = "IN"
)

type Filter struct {
	expression string
	items      FilterItems
}

func (f *Filter) Items() sqlcraft.FilterItems {
	items := sqlcraft.FilterItems{}
	for _, v := range f.items {
		items = append(items, v)
	}

	return items
}

func NewFilter(items FilterItems) Filter {
	return Filter{items: items}
}

func NewFilterItems(items ...FilterItem) Filter {
	return Filter{
		items: items,
	}
}

func NewFilterExpression(expression string) Filter {
	return Filter{expression: expression}
}

func (f *Filter) AppendItems(items ...FilterItem) {
	f.items = append(f.items, items...)
}

func (f *Filter) PushItems(items ...FilterItem) {
	f.items = append(items, f.items...)
}

func (f *Filter) ItemsLen() int {
	return len(f.items)
}

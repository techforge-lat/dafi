package dafi

import (
	"github.com/techforge-lat/sqlcraft"
	"strings"
)

var (
	AscOrder  = "ASC"
	DescOrder = "DESC"
)

type SortItem struct {
	Field string
	Order string
}

func (s SortItem) GetField() string {
	return s.Field
}

func (s SortItem) GetOrder() string {
	return s.Order
}

func NewSortItem(field, order string) SortItem {
	return SortItem{
		Field: field,
		Order: order,
	}
}

type SortItems []SortItem

type Sort struct {
	expression string
	items      SortItems
}

func (s Sort) Items() sqlcraft.SortItems {
	items := sqlcraft.SortItems{}
	for _, v := range s.items {
		items = append(items, v)
	}

	return items
}

func NewSort(items ...SortItem) Sort {
	return Sort{items: items}
}

func NewSortExpression(expression string) Sort {
	return Sort{expression: expression}
}

func BuildSortItems(expression string) SortItems {
	if expression == "" {
		return nil
	}
	var items SortItems

	sortParts := strings.Split(expression, ":")
	for _, v := range sortParts {
		order := v[len(v)-1]
		if order == '+' {
			items = append(items, SortItem{Field: v[:len(v)-1], Order: "ASC"})
			continue
		}

		if order == '-' {
			items = append(items, SortItem{Field: v[:len(v)-1], Order: "DESC"})
			continue
		}

		items = append(items, SortItem{Field: v, Order: "ASC"})
	}

	return items
}

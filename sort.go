package dafi

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	AscOrder  = "ASC"
	DescOrder = "DESC"
)

type Sort struct {
	expression string
	items      SortItems
}

func NewSort(items SortItems) Sort {
	return Sort{items: items}
}

func NewSortExpression(expression string) Sort {
	return Sort{expression: expression}
}

func (s *Sort) ReplaceAbstractNames(names map[string]string) error {
	if len(s.items) == 0 {
		s.items = buildSortItems(s.expression)
	}

	for i, v := range s.items {
		var isFound bool
		for abstractName, name := range names {
			if strings.HasPrefix(v.Field, abstractName) {
				s.items[i].Field = name
				isFound = true
				break
			}
		}

		if !isFound {
			return fmt.Errorf("sort: %w, field %s is missing", ErrInvalidFieldName, v.Field)
		}
	}

	return nil
}

func (s Sort) SQL() string {
	if len(s.items) == 0 {
		s.items = buildSortItems(s.expression)
	}

	if len(s.items) == 0 {
		return ""
	}

	builder := bytes.Buffer{}
	builder.WriteString(" ORDER BY ")
	for _, v := range s.items {
		builder.WriteString(v.Field)
		builder.WriteString(" ")
		builder.WriteString(v.Order)
		builder.WriteString(", ")
	}

	// removes the last `, `
	builder.Truncate(builder.Len() - 2)

	return builder.String()
}

func buildSortItems(expression string) SortItems {
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

type SortItem struct {
	Field string
	Order string
}

type SortItems []SortItem

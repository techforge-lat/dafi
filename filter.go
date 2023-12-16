package dafi

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidFilterFormat = errors.New("invalid filter format, format is: field operator [value] chainingKey")
	ErrInvalidOperator     = errors.New("invalid operator, must be a =, <, >, <=, >=, <>, !=, IS, ILIKE, LIKE")
	ErrInvalidFieldName    = errors.New("invalid field name")
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

func (f Filter) ItemsLen() int {
	return len(f.items)
}

func (f Filter) ReplaceAbstractNames(names map[string]string) error {
	for i, v := range f.items {
		var isFound bool
		for abstractName, name := range names {
			if v.Field == abstractName {
				f.items[i].Field = name
				isFound = true
				break
			}
		}

		if !isFound {
			return fmt.Errorf("filter: %w, field %s is missing", ErrInvalidFieldName, v.Field)
		}
	}

	return nil
}

func (f Filter) SQL() (string, []any, error) {
	if len(f.items) == 0 {
		ms, err := buildFilterItems(f.expression)
		if err != nil {
			return "", nil, err
		}

		f.items = ms
	}

	if len(f.items) == 0 {
		return "", nil, nil
	}

	builder := bytes.Buffer{}
	builder.WriteString(" WHERE ")

	args := []any{}

	var count int
	for index, item := range f.items {
		op, err := item.getOperator()
		if err != nil {
			return "", nil, err
		}

		builder.WriteString(item.Field)
		builder.WriteString(" ")
		builder.WriteString(op)

		if op == In {
			in, inArgs := buildIn(item.Value, count)
			builder.WriteString(" ")
			builder.WriteString(in)

			count += len(inArgs)
			args = append(args, inArgs...)
		} else {
			builder.WriteString(" ")
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(count + 1))
			count++
		}

		if item.ChainingKey != "" && len(f.items)-1 > index {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(item.ChainingKey))
			builder.WriteString(" ")
		}

		if op == In {
			continue
		}

		args = append(args, item.Value)
	}

	return strings.TrimSpace(builder.String()), args, nil
}

func isChainingKey(value string) bool {
	return strings.EqualFold(value, "AND") || strings.EqualFold(value, "OR")
}

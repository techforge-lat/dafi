package dafi

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidFilterFormat = errors.New("invalid filter format, must contain a name, operator and value in the form of name:operator:value with and optional chainingKey (and, or)")
	ErrInvalidOperator     = errors.New("invalid operator, must be a =, <, >, <=, >=, <>, !=, IS, ILIKE, LIKE")
)

type Filter struct {
	expression string
	items      FilterItems
	err        error
}

func NewFilter(items FilterItems) *Filter {
	return &Filter{items: items}
}

func NewFilterExpression(expression string) *Filter {
	return &Filter{expression: expression}
}

func (f *Filter) SQL() (string, []any) {
	if len(f.items) == 0 {
		f.items, f.err = buildFilterItems(f.expression)
		if f.err != nil {
			return "", nil
		}
	}

	if len(f.items) == 0 {
		return "", nil
	}

	builder := bytes.Buffer{}
	builder.WriteString(" WHERE ")

	args := []any{}

	for index, item := range f.items {
		op, err := item.getOperator()
		if err != nil {
			f.err = err
			return "", nil
		}

		builder.WriteString(item.Field)
		builder.WriteString(" ")
		builder.WriteString(op)
		builder.WriteString(" ")
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(index + 1))

		if item.ChainingKey != "" {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(item.ChainingKey))
			builder.WriteString(" ")
		}

		args = append(args, item.Value)
	}

	return strings.TrimSpace(builder.String()), args
}

func (f *Filter) Err() error {
	return f.err
}

type FilterItem struct {
	Field       string
	Operator    string
	Value       any
	ChainingKey string
}

func NewFilterItem() *FilterItem {
	return &FilterItem{}
}

func (f *FilterItem) SetField(field string) *FilterItem {
	f.Field = field
	return f
}

func (f *FilterItem) SetOperator(op string) *FilterItem {
	f.Operator = op
	return f
}

func (f *FilterItem) SetValue(value string) *FilterItem {
	f.Value = value
	return f
}

func (f *FilterItem) SetChainingKey(key string) *FilterItem {
	f.ChainingKey = key
	return f
}

func (f FilterItem) getOperator() (string, error) {
	validOperators := []string{"=", "<", ">", "<=", ">=", "<>", "!=", "is", "ilike", "like", "not_like"}

	for _, v := range validOperators {
		if strings.EqualFold(f.Operator, "not_like") {
			return "NOT LIKE", nil
		}

		if strings.EqualFold(f.Operator, "is_not") {
			return "IS NOT", nil
		}

		if strings.EqualFold(f.Operator, v) {
			return strings.ToUpper(v), nil
		}
	}

	return "", ErrInvalidOperator
}

type FilterItems []FilterItem

func buildFilterItems(expression string) (FilterItems, error) {
	if expression == "" {
		return nil, nil
	}
	var items FilterItems

	queryParts := strings.Split(expression, ";")
	for _, v := range queryParts {
		firstParts := strings.Split(v, " ")[:2]
		firstPartsLen := len(firstParts)
		if firstPartsLen != 2 {
			return nil, ErrInvalidFilterFormat
		}

		valueOpenIndex := strings.Index(v, "[")
		valueCloseIndex := strings.LastIndex(v, "]")
		item := FilterItem{
			Field:       firstParts[0],
			Operator:    firstParts[1],
			Value:       v[valueOpenIndex+1 : valueCloseIndex],
			ChainingKey: strings.TrimSpace(v[valueCloseIndex+1:]),
		}

		items = append(items, item)
	}

	return items, nil
}

func isChainingKey(value string) bool {
	return strings.EqualFold(value, "AND") || strings.EqualFold(value, "OR")
}

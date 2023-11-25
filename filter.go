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
	Different          = "<>"
	Is                 = "IS"
	IsNot              = "IS_NOT"
	ILike              = "ILIKE"
	NotILike           = "NOT_ILIKE"
	Like               = "LIKE"
	NotLike            = "NOT_LIKE"
)

type Filter struct {
	expression string
	items      FilterItems
}

func NewFilter(items FilterItems) Filter {
	return Filter{items: items}
}

func NewFilterExpression(expression string) Filter {
	return Filter{expression: expression}
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

	for index, item := range f.items {
		op, err := item.getOperator()
		if err != nil {
			return "", nil, err
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

	return strings.TrimSpace(builder.String()), args, nil
}

type FilterItem struct {
	Field       string
	Operator    string
	Value       any
	ChainingKey string
}

func (f FilterItem) getOperator() (string, error) {
	validOperators := []string{"=", "<", ">", "<=", ">=", "<>", "!=", "is", "is_not", "ilike", "not_ilike", "like", "not_like"}

	for _, v := range validOperators {
		if strings.EqualFold(f.Operator, v) {
			return strings.ReplaceAll(strings.ToUpper(v), "_", " "), nil
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

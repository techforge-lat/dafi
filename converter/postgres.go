package converter

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/techforge-lat/dafi"
)

var psqlOperatorByDafiOperator = map[dafi.FilterOperator]string{
	dafi.Equal:          "=",
	dafi.NotEqual:       "<>",
	dafi.Greater:        ">",
	dafi.GreaterOrEqual: ">=",
	dafi.Less:           "<",
	dafi.LessOrEqual:    "<=",
	dafi.Contains:       "ILIKE",
	dafi.NotContains:    "NOT ILIKE",
	dafi.Is:             "IS",
	dafi.IsNot:          "IS NOT",
	dafi.In:             "IN",
	dafi.NotIn:          "NOT IN",
}

type PsqlResult struct {
	Sql  string
	Args []any
}

type PsqlConverter struct {
	MaxPageSize uint
}

func NewPsqlConverter(maxPageSize uint) PsqlConverter {
	return PsqlConverter{
		MaxPageSize: maxPageSize,
	}
}

func (p PsqlConverter) ToSQL(criteria dafi.Criteria) (PsqlResult, error) {
	whereSql, args, err := p.BuildWhere(criteria.Filters)
	if err != nil {
		return PsqlResult{}, err
	}

	sortSql := p.BuildSort(criteria.Sorts)
	paginationSql := p.BuildPagination(criteria.Pagination)

	builder := strings.Builder{}
	builder.WriteString(whereSql)

	if sortSql != "" {
		builder.WriteString(" ")
		builder.WriteString(sortSql)
	}

	if paginationSql != "" {
		builder.WriteString(" ")
		builder.WriteString(paginationSql)
	}

	return PsqlResult{
		Sql:  strings.TrimSpace(builder.String()),
		Args: args,
	}, nil
}

func (p PsqlConverter) BuildWhere(filters dafi.Filters) (string, []any, error) {
	if filters.IsZero() {
		return "", nil, nil
	}

	builder := strings.Builder{}
	args := []any{}

	builder.WriteString("WHERE ")
	for i, filter := range filters {
		if filter.IsGroupOpen {
			if filter.GroupOpenQty == 0 {
				filter.GroupOpenQty = 1
			}

			builder.WriteString(strings.Repeat("(", filter.GroupOpenQty))
		}

		operator, ok := psqlOperatorByDafiOperator[filter.Operator]
		if !ok {
			return "", nil, errors.Join(fmt.Errorf("operator %q not found", filter.Operator), ErrInvalidOperator)
		}

		if filter.Operator == dafi.In || filter.Operator == dafi.NotIn {
			in, inArgs := p.BuildIn(filter.Value, len(args)+1)
			if in == "" {
				continue
			}

			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)

			builder.WriteString(" ")
			builder.WriteString(in)

			args = append(args, inArgs...)
		} else {
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)
			builder.WriteString(" $")
			builder.WriteString(strconv.Itoa(i + 1))

			args = append(args, filter.Value)
		}

		if i < len(filters)-1 && filter.ChainingKey == "" {
			filter.ChainingKey = dafi.And
		}

		if filter.IsGroupClose {
			if filter.GroupCloseQty == 0 {
				filter.GroupCloseQty = 1
			}

			builder.WriteString(strings.Repeat(")", filter.GroupCloseQty))
		}

		builder.WriteString(" ")
		builder.WriteString(string(filter.ChainingKey))
		builder.WriteString(" ")
	}

	return strings.TrimSpace(builder.String()), args, nil
}

func (p PsqlConverter) BuildIn(value any, index int) (string, []any) {
	if value == nil {
		return "", nil
	}

	builder := bytes.Buffer{}
	builder.WriteString("(")

	var args []any

	// uses reflection to handle different types
	valSlice := reflect.ValueOf(value)
	if valSlice.Kind() == reflect.Slice {
		if valSlice.Len() == 0 {
			return "", nil
		}

		for i := 0; i < valSlice.Len(); i++ {
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(index + i))
			builder.WriteString(", ")

			args = append(args, valSlice.Index(i).Interface())
		}

		if valSlice.Len() > 0 {
			builder.Truncate(builder.Len() - 2)
		}

		builder.WriteString(")")

		return builder.String(), args
	}

	str, ok := value.(string)
	if !ok {
		return "", nil
	}

	stringValues := strings.Split(str, ",")
	for i, v := range stringValues {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(index + i))
		builder.WriteString(", ")

		args = append(args, v)
	}

	builder.Truncate(builder.Len() - 2)
	builder.WriteString(")")

	return builder.String(), args
}

func (p PsqlConverter) BuildSort(sorts dafi.Sorts) string {
	if sorts.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("ORDER BY ")
	for i, sort := range sorts {
		builder.WriteString(string(sort.Field))

		if sort.Type != dafi.None {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(string(sort.Type)))
		}

		if i < len(sorts)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func (p PsqlConverter) BuildPagination(pagination dafi.Pagination) string {
	if !pagination.HasPageSize() {
		pagination.PageSize = uint(p.MaxPageSize)
	}

	if pagination.HasPageSize() && !pagination.HasPageNumber() {
		pagination.PageNumber = 1
	}

	if pagination.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("LIMIT ")
	builder.WriteString(strconv.Itoa(int(pagination.PageSize)))

	if pagination.HasPageNumber() {
		builder.WriteString(" OFFSET ")
		builder.WriteString(strconv.Itoa(int(pagination.PageSize * (pagination.PageNumber - 1))))
	}

	return builder.String()
}

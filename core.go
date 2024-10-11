package dafi

import (
	"regexp"
	"strconv"
	"strings"
)

type valueType string

const (
	stringType valueType = "string"
	intType    valueType = "int"
	floatType  valueType = "float"
	boolType   valueType = "bool"
)

var filterRegex = regexp.MustCompile(`(?i)^(\(+)?\s*@([a-z_]+)\s*(=|!=|>|>=|<|<=|CONTAINS|NOT_CONTAINS|IS|IS_NOT|IN|NOT_IN)\s*\[([^\]]+)\]\s*(\)+)?\s*(AND|OR)?$`)

const (
	groupOpenIndex = iota + 1
	fieldIndex
	operatorIndex
	valueIndex
	groupCloseIndex
	chainingKeyIndex
)

type Criteria struct {
	Filters    Filters
	Sorts      Sorts
	Pagination Pagination
}

func New() Criteria {
	return Criteria{}
}

func ParseFilters(expressions []string) (Filters, error) {
	if len(expressions) == 0 {
		return nil, nil
	}

	filters := Filters{}

	for _, expression := range expressions {
		if !filterRegex.MatchString(expression) {
			return nil, ErrInvalidFilterFormat
		}

		matches := filterRegex.FindStringSubmatch(expression)

		var value any

		operator := FilterOperator(matches[operatorIndex])
		valueStr := matches[valueIndex]
		value = valueStr

		if operator == In || operator == NotIn {
			valueItems, err := parseInValue(valueStr)
			if err != nil {
				return nil, err
			}
			value = valueItems
		}

		filters = append(filters, Filter{
			IsGroupOpen:   !isEmpty(matches[groupOpenIndex]),
			GroupOpenQty:  len(matches[groupOpenIndex]),
			Field:         FilterField(matches[fieldIndex]),
			Operator:      operator,
			Value:         value,
			IsGroupClose:  !isEmpty(matches[groupCloseIndex]),
			GroupCloseQty: len(matches[groupCloseIndex]),
			ChainingKey:   FilterChainingKey(matches[chainingKeyIndex]),
		})
	}

	return filters, nil
}

func ParseSorts(expression string) Sorts {
	return Sorts{}
}

func ParsePagination(expression string) Pagination {
	return Pagination{}
}

func isEmpty(v string) bool {
	return v == ""
}

func parseInValue(valueStr string) (any, error) {
	valueItems := strings.Split(valueStr, ",")
	valueType := identifyType(valueItems[0])
	missMatchValueType := false

	for i := range valueItems {
		valueItems[i] = strings.TrimSpace(valueItems[i])

		newValueType := identifyType(valueItems[i])

		if valueType != newValueType {
			missMatchValueType = true
		}

		// if we're dealing with ints and floats, we just say every item is a float
		if (valueType == intType && newValueType == floatType) || (valueType == floatType && newValueType == intType) {
			missMatchValueType = false
		}

		// we give priority to floats
		if valueType == intType && newValueType == floatType {
			valueType = floatType
		}
	}

	if len(valueItems) == 0 || missMatchValueType {
		return valueItems, nil
	}

	switch valueType {
	case intType:
		values, err := convertToIntSlice(valueItems)
		if err != nil {
			return nil, err
		}

		return values, nil
	case floatType:
		values, err := convertToFloatSlice(valueItems)
		if err != nil {
			return nil, err
		}

		return values, nil
	case boolType:
		values, err := convertToBoolSlice(valueItems)
		if err != nil {
			return nil, err
		}

		return values, nil
	}

	return valueItems, nil
}

func identifyType(value string) valueType {
	if _, err := strconv.Atoi(value); err == nil {
		return intType
	}

	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return floatType
	}

	if _, err := strconv.ParseBool(value); err == nil {
		return boolType
	}

	return stringType
}

func convertToIntSlice(values []string) ([]int, error) {
	ints := make([]int, len(values))
	for i, s := range values {
		value, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		ints[i] = value
	}

	return ints, nil
}

func convertToFloatSlice(values []string) ([]float64, error) {
	floats := make([]float64, len(values))
	for i, s := range values {
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}

		floats[i] = value
	}

	return floats, nil
}

func convertToBoolSlice(values []string) ([]bool, error) {
	bools := make([]bool, len(values))
	for i, s := range values {
		value, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}

		bools[i] = value
	}

	return bools, nil
}

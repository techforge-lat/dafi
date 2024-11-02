package dafi

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	parameterPage   = "page"
	parameterLimit  = "limit"
	parameterSort   = "sort"
	defaultChaining = And
)

type QueryParser struct {
	operators map[FilterOperator]struct{}
}

func NewQueryParser() *QueryParser {
	return &QueryParser{
		operators: map[FilterOperator]struct{}{
			Equal:          {},
			NotEqual:       {},
			Greater:        {},
			GreaterOrEqual: {},
			Less:           {},
			LessOrEqual:    {},
			Like:           {},
			In:             {},
			NotIn:          {},
			Contains:       {},
			NotContains:    {},
			Is:             {},
			IsNot:          {},
		},
	}
}

func (p *QueryParser) Parse(values url.Values) (Criteria, error) {
	criteria := Criteria{}

	for key, values := range values {
		if err := p.parseValues(key, values, &criteria); err != nil {
			return Criteria{}, err
		}
	}

	return criteria, nil
}

func (p *QueryParser) parseValues(key string, values []string, criteria *Criteria) error {
	for _, value := range values {
		if value == "" {
			continue
		}

		parts := strings.SplitN(value, ":", 3)
		if err := p.parsePart(key, parts, criteria); err != nil {
			return err
		}
	}

	return nil
}

func (p *QueryParser) parsePart(key string, parts []string, criteria *Criteria) error {
	switch {
	case p.isPaginationPart(parts):
		return p.parsePagination(parts, &criteria.Pagination)
	case p.isSortPart(parts):
		criteria.Sorts = append(criteria.Sorts, p.parseSort(key, parts))
	default:
		filter, err := p.parseFilter(key, parts)
		if err != nil {
			return err
		}

		criteria.Filters = append(criteria.Filters, filter)
	}

	return nil
}

func (p *QueryParser) isPaginationPart(parts []string) bool {
	return len(parts) == 2 && (parts[0] == parameterPage || parts[0] == parameterLimit)
}

func (p *QueryParser) isSortPart(parts []string) bool {
	return len(parts) == 2 && parts[0] == parameterSort
}

func (p *QueryParser) parsePagination(parts []string, pagination *Pagination) error {
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	switch parts[0] {
	case parameterPage:
		pagination.PageNumber = uint(value)
	case parameterLimit:
		pagination.PageSize = uint(value)
	}

	return nil
}

func (p *QueryParser) parseSort(field string, parts []string) Sort {
	return Sort{
		Field: SortBy(field),
		Type:  SortType(strings.ToUpper(parts[1])),
	}
}

func (p *QueryParser) parseFilter(field string, parts []string) (Filter, error) {
	operator := p.determineOperator(parts[0])
	chainingKey := p.determineChainingKey(parts)

	var value any = parts[1]
	if operator == In || operator == NotIn {
		value = strings.Split(parts[1], ",")
	}

	return Filter{
		Field:       FilterField(field),
		Operator:    operator,
		Value:       value,
		ChainingKey: chainingKey,
	}, nil
}

func (p *QueryParser) determineOperator(op string) FilterOperator {
	operator := FilterOperator(op)
	if _, ok := p.operators[operator]; !ok {
		return Equal
	}

	return operator
}

func (p *QueryParser) determineChainingKey(parts []string) FilterChainingKey {
	if len(parts) == 3 {
		return FilterChainingKey(strings.ToUpper(parts[2]))
	}

	return defaultChaining
}

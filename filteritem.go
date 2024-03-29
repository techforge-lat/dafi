package dafi

import (
	"strings"
)

type FilterItem struct {
	Field       string
	Operator    string
	Value       any
	ChainingKey string
	GroupOpen   string
	GroupClose  string
}

func NewFilterItem(field string, operator string, value any, chainingKey ...string) FilterItem {
	f := FilterItem{
		Field:    field,
		Operator: operator,
		Value:    value,
	}

	if len(chainingKey) > 0 {
		f.ChainingKey = chainingKey[0]
	}

	return f
}

func (f FilterItem) GetField() string {
	return f.Field
}

func (f FilterItem) GetOperator() string {
	return string(f.Operator)
}

func (f FilterItem) GetValue() any {
	return f.Value
}

func (f FilterItem) GetChainingKey() string {
	return f.ChainingKey
}

func (f FilterItem) GetGroupOpen() string {
	return f.GroupOpen
}

func (f FilterItem) GetGroupClose() string {
	return f.GroupClose
}

func (f FilterItem) HasGroupOpen() bool {
	return f.GroupOpen != ""
}

func (f FilterItem) HasGroupClose() bool {
	return f.GroupClose != ""
}

func (f FilterItem) GetParsedOperator() (string, error) {
	validOperators := []string{"=", "<", ">", "<=", ">=", "<>", "!=", "is", "is_not", "ilike", "not_ilike", "like", "not_like", "in"}

	for _, v := range validOperators {
		if strings.EqualFold(f.Operator, v) {
			return strings.ReplaceAll(strings.ToUpper(v), "_", " "), nil
		}
	}

	return "", ErrInvalidOperator
}

type FilterItems []FilterItem

func BuildFilterItems(expression string) (FilterItems, error) {
	if expression == "" {
		return nil, nil
	}
	var items FilterItems

	queryParts := strings.Split(expression, ":")
	for _, v := range queryParts {
		firstParts := strings.Split(v, " ")
		firstPartsLen := len(firstParts)

		if firstPartsLen == 1 && firstParts[0] == "(" {
			item := FilterItem{
				GroupOpen: "(",
			}

			items = append(items, item)
			continue
		}

		if firstPartsLen == 2 && firstParts[1] == "(" {
			item := FilterItem{
				GroupOpen:   "(",
				ChainingKey: firstParts[0],
			}

			items = append(items, item)
			continue
		}

		if firstPartsLen == 1 && firstParts[0] == ")" {
			item := FilterItem{
				GroupClose: ")",
			}

			items = append(items, item)
			continue
		}

		if firstPartsLen == 2 && firstParts[0] == ")" {
			item := FilterItem{
				GroupClose:  ")",
				ChainingKey: firstParts[1],
			}

			items = append(items, item)
			continue
		}

		firstParts = firstParts[:2]
		firstPartsLen = len(firstParts)
		if firstPartsLen != 2 {
			return nil, ErrInvalidFilterFormat
		}

		valueOpenIndex := strings.Index(v, "[")
		if valueOpenIndex == -1 {
			return nil, ErrInvalidFilterFormat
		}

		valueCloseIndex := strings.LastIndex(v, "]")
		if valueCloseIndex == -1 {
			return nil, ErrInvalidFilterFormat
		}

		item := FilterItem{
			Field:       firstParts[0],
			Operator:    strings.ToUpper(firstParts[1]),
			Value:       v[valueOpenIndex+1 : valueCloseIndex],
			ChainingKey: strings.TrimSpace(v[valueCloseIndex+1:]),
		}

		items = append(items, item)
	}

	return items, nil
}

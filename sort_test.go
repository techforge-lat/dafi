package dafi

import (
	"reflect"
	"testing"
)

func TestBuildSortItems(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name   string
		fields fields
		want   SortItems
	}{
		{
			name: "",
			fields: fields{
				expression: "name+:age-",
			},
			want: SortItems{
				{Field: "name", Order: "ASC"},
				{Field: "age", Order: "DESC"},
			},
		},
		{
			name: "",
			fields: fields{
				expression: "name+:age",
			},
			want: SortItems{
				{Field: "name", Order: "ASC"},
				{Field: "age", Order: "ASC"},
			},
		},
		{
			name: "",
			fields: fields{
				expression: "name:age",
			},
			want: SortItems{
				{Field: "name", Order: "ASC"},
				{Field: "age", Order: "ASC"},
			},
		},
		{
			name: "",
			fields: fields{
				expression: "name",
			},
			want: SortItems{
				{Field: "name", Order: "ASC"},
			},
		},
		{
			name: "",
			fields: fields{
				expression: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildSortItems(tt.fields.expression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildSortItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

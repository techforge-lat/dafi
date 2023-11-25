package dafi

import "testing"

func TestSort_SQL(t *testing.T) {
	type fields struct {
		expression string
		items      SortItems
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				expression: "name+;age-",
			},
			want: " ORDER BY name ASC, age DESC",
		},
		{
			name: "",
			fields: fields{
				expression: "name+;age",
			},
			want: " ORDER BY name ASC, age ASC",
		},
		{
			name: "",
			fields: fields{
				expression: "name;age",
			},
			want: " ORDER BY name ASC, age ASC",
		},
		{
			name: "",
			fields: fields{
				expression: "name",
			},
			want: " ORDER BY name ASC",
		},
		{
			name: "",
			fields: fields{
				expression: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sort{
				expression: tt.fields.expression,
				items:      tt.fields.items,
			}
			if got := s.SQL(); got != tt.want {
				t.Errorf("Sort.SQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

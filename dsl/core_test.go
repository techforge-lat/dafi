package dsl

import (
	"reflect"
	"testing"

	"github.com/techforge-lat/dafi"
)

func TestDSL_ParseFilters(t *testing.T) {
	type args struct {
		expressions []string
	}
	tests := []struct {
		name    string
		d       DSL
		args    args
		want    dafi.Filters
		wantErr bool
	}{
		{
			name: "Single equality filter on email",
			d:    New(),
			args: args{
				expressions: []string{"@email = [hernan_rm@outlook.es]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "email",
					Operator:     dafi.Equal,
					Value:        "hernan_rm@outlook.es",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not-equal filter on email",
			d:    New(),
			args: args{
				expressions: []string{"@email != [hernan_rm@outlook.es]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "email",
					Operator:     dafi.NotEqual,
					Value:        "hernan_rm@outlook.es",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single greater than filter on age",
			d:    New(),
			args: args{
				expressions: []string{"@age > [30]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     dafi.Greater,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single greater than or equal filter on age",
			d:    New(),
			args: args{
				expressions: []string{"@age >= [30]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     dafi.GreaterOrEqual,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single less than filter on age",
			d:    New(),
			args: args{
				expressions: []string{"@age < [30]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     dafi.Less,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single less than or equal filter on age",
			d:    New(),
			args: args{
				expressions: []string{"@age <= [30]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     dafi.LessOrEqual,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single contains filter on name",
			d:    New(),
			args: args{
				expressions: []string{"@name CONTAINS [John]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     dafi.Contains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not contains filter on name",
			d:    New(),
			args: args{
				expressions: []string{"@name NOT_CONTAINS [John]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     dafi.NotContains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single is filter on active status",
			d:    New(),
			args: args{
				expressions: []string{"@active IS [true]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "active",
					Operator:     dafi.Is,
					Value:        "true",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single is not filter on active status",
			d:    New(),
			args: args{
				expressions: []string{"@active IS_NOT [false]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "active",
					Operator:     dafi.IsNot,
					Value:        "false",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single in filter on user IDs",
			d:    New(),
			args: args{
				expressions: []string{"@user_id IN [1, 2, 3]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     dafi.In,
					Value:        []int{1, 2, 3},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single in filter on user nicknames with numbers",
			d:    New(),
			args: args{
				expressions: []string{"@nickname IN [2024, yanelly, hernan]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "nickname",
					Operator:     dafi.In,
					Value:        []string{"2024", "yanelly", "hernan"},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on user IDs",
			d:    New(),
			args: args{
				expressions: []string{"@user_id NOT_IN [4, 5, 6]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     dafi.NotIn,
					Value:        []int{4, 5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on prices in",
			d:    New(),
			args: args{
				expressions: []string{"@user_id IN [4.4, 5, 6]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     dafi.In,
					Value:        []float64{4.4, 5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on prices in starting with int but having floats after",
			d:    New(),
			args: args{
				expressions: []string{"@user_id IN [4, 5.5, 6]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     dafi.In,
					Value:        []float64{4, 5.5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid filter format with missing operator",
			d:    New(),
			args: args{
				expressions: []string{"@age 30"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid filter format with unsupported operator",
			d:    New(),
			args: args{
				expressions: []string{"@age => 30"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty filter expression",
			d:    New(),
			args: args{
				expressions: []string{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Multiple filters",
			d:    New(),
			args: args{
				expressions: []string{"@age > [30] AND", "@name CONTAINS [John]"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     dafi.Greater,
					Value:        "30",
					ChainingKey:  dafi.And,
					IsGroupClose: false,
				},
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     dafi.Contains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple filters with groups",
			d:    New(),
			args: args{
				expressions: []string{"(( @age > [30] ) AND", "( @name CONTAINS [John] ))"},
			},
			want: dafi.Filters{
				{
					IsGroupOpen:   true,
					GroupOpenQty:  2,
					Field:         "age",
					Operator:      dafi.Greater,
					Value:         "30",
					ChainingKey:   dafi.And,
					IsGroupClose:  true,
					GroupCloseQty: 1,
				},
				{
					IsGroupOpen:   true,
					GroupOpenQty:  1,
					Field:         "name",
					Operator:      dafi.Contains,
					Value:         "John",
					ChainingKey:   "",
					IsGroupClose:  true,
					GroupCloseQty: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.ParseFilters(tt.args.expressions)
			if (err != nil) != tt.wantErr {
				t.Errorf("DSL.ParseFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DSL.ParseFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

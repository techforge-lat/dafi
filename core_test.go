package dafi

import (
	"reflect"
	"testing"
)

func TestDSL_ParseFilters(t *testing.T) {
	type args struct {
		expressions []string
	}
	tests := []struct {
		name    string
		args    args
		want    Filters
		wantErr bool
	}{
		{
			name: "Single equality filter on email",
			args: args{
				expressions: []string{"@email = [hernan_rm@outlook.es]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "email",
					Operator:     Equal,
					Value:        "hernan_rm@outlook.es",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not-equal filter on email",
			args: args{
				expressions: []string{"@email != [hernan_rm@outlook.es]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "email",
					Operator:     NotEqual,
					Value:        "hernan_rm@outlook.es",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single greater than filter on age",
			args: args{
				expressions: []string{"@age > [30]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     Greater,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single greater than or equal filter on age",
			args: args{
				expressions: []string{"@age >= [30]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     GreaterOrEqual,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single less than filter on age",
			args: args{
				expressions: []string{"@age < [30]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     Less,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single less than or equal filter on age",
			args: args{
				expressions: []string{"@age <= [30]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     LessOrEqual,
					Value:        "30",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single contains filter on name",
			args: args{
				expressions: []string{"@name CONTAINS [John]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     Contains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not contains filter on name",
			args: args{
				expressions: []string{"@name NOT_CONTAINS [John]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     NotContains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single is filter on active status",
			args: args{
				expressions: []string{"@active IS [true]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "active",
					Operator:     Is,
					Value:        "true",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single is not filter on active status",
			args: args{
				expressions: []string{"@active IS_NOT [false]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "active",
					Operator:     IsNot,
					Value:        "false",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single in filter on user IDs",
			args: args{
				expressions: []string{"@user_id IN [1, 2, 3]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     In,
					Value:        []int{1, 2, 3},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single in filter on user nicknames with numbers",
			args: args{
				expressions: []string{"@nickname IN [2024, yanelly, hernan]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "nickname",
					Operator:     In,
					Value:        []string{"2024", "yanelly", "hernan"},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on user IDs",
			args: args{
				expressions: []string{"@user_id NOT_IN [4, 5, 6]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     NotIn,
					Value:        []int{4, 5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on prices in",
			args: args{
				expressions: []string{"@user_id IN [4.4, 5, 6]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     In,
					Value:        []float64{4.4, 5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Single not in filter on prices in starting with int but having floats after",
			args: args{
				expressions: []string{"@user_id IN [4, 5.5, 6]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "user_id",
					Operator:     In,
					Value:        []float64{4, 5.5, 6},
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid filter format with missing operator",
			args: args{
				expressions: []string{"@age 30"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid filter format with unsupported operator",
			args: args{
				expressions: []string{"@age => 30"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty filter expression",
			args: args{
				expressions: []string{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Multiple filters",
			args: args{
				expressions: []string{"@age > [30] AND", "@name CONTAINS [John]"},
			},
			want: Filters{
				{
					IsGroupOpen:  false,
					Field:        "age",
					Operator:     Greater,
					Value:        "30",
					ChainingKey:  And,
					IsGroupClose: false,
				},
				{
					IsGroupOpen:  false,
					Field:        "name",
					Operator:     Contains,
					Value:        "John",
					ChainingKey:  "",
					IsGroupClose: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple filters with groups",
			args: args{
				expressions: []string{"(( @age > [30] ) AND", "( @name CONTAINS [John] ))"},
			},
			want: Filters{
				{
					IsGroupOpen:   true,
					GroupOpenQty:  2,
					Field:         "age",
					Operator:      Greater,
					Value:         "30",
					ChainingKey:   And,
					IsGroupClose:  true,
					GroupCloseQty: 1,
				},
				{
					IsGroupOpen:   true,
					GroupOpenQty:  1,
					Field:         "name",
					Operator:      Contains,
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
			got, err := ParseFilters(tt.args.expressions)
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

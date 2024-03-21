package dafi

import (
	"reflect"
	"testing"
)

func TestBuildFilterItemsFromExpression(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		args    args
		want    FilterItems
		wantErr bool
	}{
		{
			name: "username and email",
			args: args{
				expression: "username = [john doe] AND:email LIKE [%example.com]",
			},
			want: FilterItems{
				{
					Field:       "username",
					Operator:    Equal,
					Value:       "john doe",
					ChainingKey: And,
				},
				{
					Field:    "email",
					Operator: Like,
					Value:    "%example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "username and is_admin",
			args: args{
				expression: "username = [admin] OR:is_admin = [true]",
			},
			want: FilterItems{
				{
					Field:       "username",
					Operator:    Equal,
					Value:       "admin",
					ChainingKey: Or,
				},
				{
					Field:    "is_admin",
					Operator: Equal,
					Value:    "true",
				},
			},
			wantErr: false,
		},
		{
			name: "email not like",
			args: args{
				expression: "email NOT_LIKE [%test%]",
			},
			want: FilterItems{
				{
					Field:    "email",
					Operator: NotLike,
					Value:    "%test%",
				},
			},
			wantErr: false,
		},
		{
			name: "username and email or is_admin",
			args: args{
				expression: "username = [john doe] AND:email LIKE [%example.com] OR:is_admin = [true]",
			},
			want: FilterItems{
				{
					Field:       "username",
					Operator:    Equal,
					Value:       "john doe",
					ChainingKey: And,
				},
				{
					Field:       "email",
					Operator:    Like,
					Value:       "%example.com",
					ChainingKey: Or,
				},
				{
					Field:    "is_admin",
					Operator: Equal,
					Value:    "true",
				},
			},
			wantErr: false,
		},
		{
			name: "username or email or is_admin",
			args: args{
				expression: "username = [john doe] OR:email LIKE [%example.com] OR:is_admin = [true]",
			},
			want: FilterItems{
				{
					Field:       "username",
					Operator:    Equal,
					Value:       "john doe",
					ChainingKey: Or,
				},
				{
					Field:       "email",
					Operator:    Like,
					Value:       "%example.com",
					ChainingKey: Or,
				},
				{
					Field:    "is_admin",
					Operator: Equal,
					Value:    "true",
				},
			},
			wantErr: false,
		},
		{
			name: "username in",
			args: args{
				expression: "username in [hernan,fer,david,erick]",
			},
			want: FilterItems{
				{
					Field:    "username",
					Operator: In,
					Value:    "hernan,fer,david,erick",
				},
			},
			wantErr: false,
		},
		{
			name: "email equal and slug in",
			args: args{
				expression: "domain = [facebook.com] AND:slug IN [facebook,apple,microsoft]",
			},
			want: FilterItems{
				{
					Field:       "domain",
					Operator:    Equal,
					Value:       "facebook.com",
					ChainingKey: And,
				},
				{
					Field:    "slug",
					Operator: In,
					Value:    "facebook,apple,microsoft",
				},
			},
			wantErr: false,
		},
		{
			name: "domain equal and slug in with groups",
			args: args{
				expression: "domain = [facebook.com] AND:(:slug IN [facebook,apple,microsoft] OR:id = [123]:)",
			},
			want: FilterItems{
				{
					Field:       "domain",
					Operator:    Equal,
					Value:       "facebook.com",
					ChainingKey: And,
				},
				{
					GroupOpen: "(",
				},
				{
					Field:       "slug",
					Operator:    In,
					Value:       "facebook,apple,microsoft",
					ChainingKey: Or,
				},
				{
					Field:    "id",
					Operator: Equal,
					Value:    "123",
				},
				{
					GroupClose: ")",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildFilterItems(tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildFilterItemsFromExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildFilterItemsFromExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

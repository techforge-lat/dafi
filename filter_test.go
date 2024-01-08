package dafi

import (
	"reflect"
	"testing"
)

func TestFilter_SQL(t *testing.T) {
	type fields struct {
		expression string
		err        error
	}
	tests := []struct {
		name      string
		fields    fields
		wantWhere string
		wantArgs  []any
		wantErr   bool
	}{
		{
			name: "username and email",
			fields: fields{
				expression: "username = [john doe] AND:email LIKE [%example.com]",
			},
			wantWhere: "WHERE username = $1 AND email LIKE $2",
			wantArgs:  []any{"john doe", "%example.com"},
		},
		{
			name: "is_admin true",
			fields: fields{
				expression: "is_admin = [true]",
			},
			wantWhere: "WHERE is_admin = $1",
			wantArgs:  []any{"true"},
		},
		{
			name: "username and is_admin",
			fields: fields{
				expression: "username = [admin] OR:is_admin = [true]",
			},
			wantWhere: "WHERE username = $1 OR is_admin = $2",
			wantArgs:  []any{"admin", "true"},
		},
		{
			name: "email not like",
			fields: fields{
				expression: "email NOT_LIKE [%test%]",
			},
			wantWhere: "WHERE email NOT LIKE $1",
			wantArgs:  []any{"%test%"},
		},
		{
			name: "username and (email or is_admin)",
			fields: fields{
				expression: "username = [john doe] AND:email LIKE [%example.com] OR:is_admin = [true]",
			},
			wantWhere: "WHERE username = $1 AND email LIKE $2 OR is_admin = $3",
			wantArgs:  []any{"john doe", "%example.com", "true"},
		},
		{
			name: "username or (email and is_admin)",
			fields: fields{
				expression: "username = [admin] OR:email LIKE [%example.com] AND:is_admin = [true]",
			},
			wantWhere: "WHERE username = $1 OR email LIKE $2 AND is_admin = $3",
			wantArgs:  []any{"admin", "%example.com", "true"},
		},
		{
			name: "empty",
			fields: fields{
				expression: "",
			},
			wantWhere: "",
			wantArgs:  nil,
		},
		{
			name: "username or (email and is_admin)",
			fields: fields{
				expression: "username ilike [%admin%]",
			},
			wantWhere: "WHERE username ILIKE $1",
			wantArgs:  []any{"%admin%"},
		},
		{
			name: "username in",
			fields: fields{
				expression: "username in [hernan,fer,david,erick]",
			},
			wantWhere: "WHERE username IN ($1, $2, $3, $4)",
			wantArgs:  []any{"hernan", "fer", "david", "erick"},
		},
		{
			name: "username in",
			fields: fields{
				expression: "username in [1,2,3]",
			},
			wantWhere: "WHERE username IN ($1, $2, $3)",
			wantArgs:  []any{"1", "2", "3"},
		},
		{
			name: "domain equal and slug in",
			fields: fields{
				expression: "domain = [facebook.com] AND:slug IN [facebook,apple,microsoft]",
			},
			wantWhere: "WHERE domain = $1 AND slug IN ($2, $3, $4)",
			wantArgs:  []any{"facebook.com", "facebook", "apple", "microsoft"},
		},
		{
			name: "domain equal and slug in with groups",
			fields: fields{
				expression: "domain = [facebook.com] AND:(:slug IN [facebook,apple,microsoft] OR:id = [123]:)",
			},
			wantWhere: "WHERE domain = $1 AND (slug IN ($2, $3, $4) OR id = $5)",
			wantArgs:  []any{"facebook.com", "facebook", "apple", "microsoft", "123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				expression: tt.fields.expression,
			}
			gotWhere, gotArgs, err := f.SQL()
			if gotWhere != tt.wantWhere {
				t.Errorf("Filter.SQL() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Filter.SQL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Filter.SQL() gotErr = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

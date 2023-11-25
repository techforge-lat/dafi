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
	}{
		{
			name: "username and email",
			fields: fields{
				expression: "username = [john doe] AND;email LIKE [%example.com]",
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
				expression: "username = [admin] OR;is_admin = [true]",
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
				expression: "username = [john doe] AND;email LIKE [%example.com] OR;is_admin = [true]",
			},
			wantWhere: "WHERE username = $1 AND email LIKE $2 OR is_admin = $3",
			wantArgs:  []any{"john doe", "%example.com", "true"},
		},
		{
			name: "username or (email and is_admin)",
			fields: fields{
				expression: "username = [admin] OR;email LIKE [%example.com] AND;is_admin = [true]",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				expression: tt.fields.expression,
				err:        tt.fields.err,
			}
			got, got1 := f.SQL()
			if got != tt.wantWhere {
				t.Errorf("Filter.SQL() gotWhere = %v, want %v", got, tt.wantWhere)
			}
			if !reflect.DeepEqual(got1, tt.wantArgs) {
				t.Errorf("Filter.SQL() gotArgs = %v, want %v", got1, tt.wantArgs)
			}
		})
	}
}

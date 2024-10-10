package converter

import (
	"reflect"
	"testing"

	"github.com/techforge-lat/dafi"
)

func TestPsqlConverter_Convert(t *testing.T) {
	psqlConverter := NewPsqlConverter(0)

	type args struct {
		criteria dafi.Criteria
	}
	tests := []struct {
		name    string
		p       PsqlConverter
		args    args
		want    PsqlResult
		wantErr bool
	}{
		{
			name: "one filter",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "email",
							Operator: dafi.Equal,
							Value:    "hernan_rm@outlook.es",
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE email = $1",
				Args: []any{"hernan_rm@outlook.es"},
			},
			wantErr: false,
		},
		{
			name: "and chaining key",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "email",
							Operator: dafi.Equal,
							Value:    "hernan_rm@outlook.es",
						},
						dafi.Filter{
							Field:    "nickname",
							Operator: dafi.Equal,
							Value:    "hernanreyes",
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE email = $1 AND nickname = $2",
				Args: []any{"hernan_rm@outlook.es", "hernanreyes"},
			},
			wantErr: false,
		},
		{
			name: "or chaining key",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:       "email",
							Operator:    dafi.Equal,
							Value:       "hernan_rm@outlook.es",
							ChainingKey: dafi.Or,
						},
						dafi.Filter{
							Field:    "nickname",
							Operator: dafi.Equal,
							Value:    "hernanreyes",
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE email = $1 OR nickname = $2",
				Args: []any{"hernan_rm@outlook.es", "hernanreyes"},
			},
			wantErr: false,
		},
		{
			name: "one condition group",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							IsGroupOpen: true,
							Field:       "email",
							Operator:    dafi.Equal,
							Value:       "hernan_rm@outlook.es",
							ChainingKey: dafi.Or,
						},
						dafi.Filter{
							Field:        "nickname",
							Operator:     dafi.Equal,
							Value:        "hernanreyes",
							IsGroupClose: true,
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE (email = $1 OR nickname = $2)",
				Args: []any{"hernan_rm@outlook.es", "hernanreyes"},
			},
			wantErr: false,
		},
		{
			name: "two conditions group",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							IsGroupOpen: true,
							Field:       "email",
							Operator:    dafi.Equal,
							Value:       "hernan_rm@outlook.es",
							ChainingKey: dafi.Or,
						},
						dafi.Filter{
							Field:        "nickname",
							Operator:     dafi.Equal,
							Value:        "hernanreyes",
							IsGroupClose: true,
						},
						dafi.Filter{
							IsGroupOpen: true,
							Field:       "phone_number",
							Operator:    dafi.Equal,
							Value:       "12345679",
							ChainingKey: dafi.Or,
						},
						dafi.Filter{
							Field:        "full_name",
							Operator:     dafi.Contains,
							Value:        "Hernan Reyes",
							IsGroupClose: true,
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE (email = $1 OR nickname = $2) AND (phone_number = $3 OR full_name ILIKE $4)",
				Args: []any{"hernan_rm@outlook.es", "hernanreyes", "12345679", "Hernan Reyes"},
			},
			wantErr: false,
		},
		{
			name: "two conditions group with multiple opening and clsing parenthesis",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							IsGroupOpen:  true,
							GroupOpenQty: 2,
							Field:        "email",
							Operator:     dafi.Equal,
							Value:        "hernan_rm@outlook.es",
							ChainingKey:  dafi.Or,
						},
						dafi.Filter{
							Field:        "nickname",
							Operator:     dafi.Equal,
							Value:        "hernanreyes",
							IsGroupClose: true,
							GroupOpenQty: 1,
						},
						dafi.Filter{
							IsGroupOpen:  true,
							GroupOpenQty: 1,
							Field:        "phone_number",
							Operator:     dafi.Equal,
							Value:        "12345679",
							ChainingKey:  dafi.Or,
						},
						dafi.Filter{
							Field:         "full_name",
							Operator:      dafi.Contains,
							Value:         "Hernan Reyes",
							IsGroupClose:  true,
							GroupCloseQty: 2,
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE ((email = $1 OR nickname = $2) AND (phone_number = $3 OR full_name ILIKE $4))",
				Args: []any{"hernan_rm@outlook.es", "hernanreyes", "12345679", "Hernan Reyes"},
			},
			wantErr: false,
		},
		{
			name: "in operator",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "id",
							Operator: dafi.In,
							Value:    []uint{1, 2, 3},
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE id IN ($1, $2, $3)",
				Args: []any{uint(1), uint(2), uint(3)},
			},
			wantErr: false,
		},
		{
			name: "not in operator",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "id",
							Operator: dafi.NotIn,
							Value:    []uint{1, 2, 3},
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE id NOT IN ($1, $2, $3)",
				Args: []any{uint(1), uint(2), uint(3)},
			},
			wantErr: false,
		},
		{
			name: "in operator with float",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: dafi.In,
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE price IN ($1, $2, $3)",
				Args: []any{1.1, 2.2, 3.3},
			},
			wantErr: false,
		},
		{
			name: "pagination",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: dafi.In,
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{},
					Pagination: dafi.Pagination{
						PageNumber: 1,
						PageSize:   10,
					},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE price IN ($1, $2, $3) LIMIT 10 OFFSET 0",
				Args: []any{1.1, 2.2, 3.3},
			},
			wantErr: false,
		},
		{
			name: "default pagination",
			p:    NewPsqlConverter(20),
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: dafi.In,
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE price IN ($1, $2, $3) LIMIT 20 OFFSET 0",
				Args: []any{1.1, 2.2, 3.3},
			},
			wantErr: false,
		},
		{
			name: "sort by",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: dafi.In,
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{
						dafi.Sort{
							Field: "created_at",
							Type:  dafi.Asc,
						},
					},
					Pagination: dafi.Pagination{
						PageNumber: 1,
						PageSize:   10,
					},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE price IN ($1, $2, $3) ORDER BY created_at ASC LIMIT 10 OFFSET 0",
				Args: []any{1.1, 2.2, 3.3},
			},
			wantErr: false,
		},
		{
			name: "sort by without type",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: dafi.In,
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{
						dafi.Sort{
							Field: "created_at",
						},
					},
					Pagination: dafi.Pagination{
						PageNumber: 1,
						PageSize:   10,
					},
				},
			},
			want: PsqlResult{
				Sql:  "WHERE price IN ($1, $2, $3) ORDER BY created_at LIMIT 10 OFFSET 0",
				Args: []any{1.1, 2.2, 3.3},
			},
			wantErr: false,
		},
		{
			name: "invalid operator",
			p:    psqlConverter,
			args: args{
				criteria: dafi.Criteria{
					Filters: dafi.Filters{
						dafi.Filter{
							Field:    "price",
							Operator: "invalid",
							Value:    []float64{1.1, 2.2, 3.3},
						},
					},
					Sorts: dafi.Sorts{
						dafi.Sort{
							Field: "created_at",
						},
					},
					Pagination: dafi.Pagination{
						PageNumber: 1,
						PageSize:   10,
					},
				},
			},
			want:    PsqlResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPsqlConverter(tt.p.MaxPageSize)
			got, err := p.ToSQL(tt.args.criteria)
			if (err != nil) != tt.wantErr {
				t.Errorf("PsqlConverter.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PsqlConverter.Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

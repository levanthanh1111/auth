package base

import (
	"testing"
)

func TestListParams_IsValid(t *testing.T) {
	type fields struct {
		Limit              int
		Offset             int
		Filters            []*Filter
		acceptedFilterKeys []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"success",
			fields{
				10,
				0,
				[]*Filter{
					{
						"name",
						OperatorEqual,
						"dung",
					},
				},
				[]string{"name", "age"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := &ListParams{
				Limit:              tt.fields.Limit,
				Offset:             tt.fields.Offset,
				Filters:            tt.fields.Filters,
				acceptedFilterKeys: tt.fields.acceptedFilterKeys,
			}
			if err := lp.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("ListParams.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

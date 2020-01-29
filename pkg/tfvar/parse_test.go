package tfvar

import (
	"testing"

	"github.com/hashicorp/terraform/configs"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestParseValues(t *testing.T) {
	type args struct {
		from map[string]UnparsedVariableValue
		vars []Variable
	}

	tests := []struct {
		name      string
		args      args
		want      []Variable
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				from: map[string]UnparsedVariableValue{
					"a": unparsedVariableValueString{str: "val", name: "a"},
				},
				vars: []Variable{
					{Name: "a", parsingMode: configs.VariableParseLiteral},
					{Name: "b", parsingMode: configs.VariableParseLiteral},
				},
			},
			want: []Variable{
				{Name: "a", Value: cty.StringVal("val"), parsingMode: configs.VariableParseLiteral},
				{Name: "b", parsingMode: configs.VariableParseLiteral},
			},
			assertion: assert.NoError,
		},
		{
			name: "failed parsing mode",
			args: args{
				from: map[string]UnparsedVariableValue{
					"a": unparsedVariableValueString{str: "val", name: "a"},
				},
				vars: []Variable{
					{Name: "a", parsingMode: configs.VariableParseHCL},
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValues(tt.args.from, tt.args.vars)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

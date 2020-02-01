package tfvar

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestCollectFromEnvVars(t *testing.T) {
	require.NoError(t, os.Setenv("MY_VAR", "my-value"))
	require.NoError(t, os.Setenv("TF_VAR_availability_zone_names", `'["us-west-1a"]'`))

	actual := make(map[string]UnparsedVariableValue)
	collectFromEnvVars(actual)

	expected := map[string]UnparsedVariableValue{
		"availability_zone_names": unparsedVariableValueString{
			str:  `'["us-west-1a"]'`,
			name: "availability_zone_names",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestCollectFromString(t *testing.T) {
	type args struct {
		raw string
		to  map[string]UnparsedVariableValue
	}

	tests := []struct {
		name      string
		args      args
		want      map[string]UnparsedVariableValue
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				raw: "a=val_a",
				to:  map[string]UnparsedVariableValue{},
			},
			want: map[string]UnparsedVariableValue{
				"a": unparsedVariableValueString{
					str:  `val_a`,
					name: "a",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "no equal",
			args: args{
				raw: "a",
				to:  map[string]UnparsedVariableValue{},
			},
			want:      map[string]UnparsedVariableValue{},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, collectFromString(tt.args.raw, tt.args.to))
			assert.Equal(t, tt.want, tt.args.to)
		})
	}
}

func TestCollectFromFile(t *testing.T) {
	type args struct {
		filename string
		to       map[string]UnparsedVariableValue
	}

	tests := []struct {
		name      string
		args      args
		want      []cty.Value
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "not found",
			args: args{
				filename: "unknown.tfvars",
				to:       map[string]UnparsedVariableValue{},
			},
			want:      []cty.Value{},
			assertion: assert.Error,
		},
		{
			name: "normal tfvars file",
			args: args{
				filename: "testdata/normal.tfvars",
				to:       map[string]UnparsedVariableValue{},
			},
			want: []cty.Value{
				cty.StringVal("<RESOURCE_PREFIX>"),
			},
			assertion: assert.NoError,
		},
		{
			name: "normal json file",
			args: args{
				filename: "testdata/normal.tfvars.json",
				to:       map[string]UnparsedVariableValue{},
			},
			want: []cty.Value{
				cty.StringVal("hello"),
			},
			assertion: assert.NoError,
		},
		{
			name: "bad tfvars file",
			args: args{
				filename: "testdata/bad.tfvars",
				to:       map[string]UnparsedVariableValue{},
			},
			want:      []cty.Value{},
			assertion: assert.Error,
		},
		{
			name: "bad json file",
			args: args{
				filename: "testdata/bad.tfvars.json",
				to:       map[string]UnparsedVariableValue{},
			},
			want:      []cty.Value{},
			assertion: assert.Error,
		},
		{
			name: "not tfvars file",
			args: args{
				filename: "testdata/normal/main.tf",
				to:       map[string]UnparsedVariableValue{},
			},
			want:      []cty.Value{},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, collectFromFile(tt.args.filename, tt.args.to))

			actual := make([]cty.Value, 0, len(tt.args.to))
			for _, v := range tt.args.to {
				val, err := v.ParseVariableValue('x')
				require.NoError(t, err)
				actual = append(actual, val)

			}
			assert.Equal(t, tt.want, actual)
		})
	}
}

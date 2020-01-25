package tfvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestLoad(t *testing.T) {
	type args struct {
		rootDir string
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
				rootDir: "./testdata/normal",
			},
			want: []Variable{
				{Name: "resource_name"},
				{Name: "instance_name", Value: cty.StringVal("my-instance")},
			},
			assertion: assert.NoError,
		},
		{
			name: "bad",
			args: args{
				rootDir: "./testdata/bad",
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.rootDir)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

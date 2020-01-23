package tfvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type args struct {
		rootDir string
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				rootDir: "./testdata/normal",
			},
			want:      []string{"resource_name"},
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

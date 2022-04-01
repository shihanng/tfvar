package tfvar

import (
	"bytes"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/configs"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestLoad(t *testing.T) {
	type args struct {
		dir string
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
				dir: "./testdata/normal",
			},
			want: []Variable{
				{Name: "resource_name", parsingMode: configs.VariableParseLiteral},
				{Name: "instance_name", Value: cty.StringVal("my-instance"), parsingMode: configs.VariableParseLiteral},
			},
			assertion: assert.NoError,
		},
		{
			name: "bad",
			args: args{
				dir: "./testdata/bad",
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.dir)
			tt.assertion(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestWriteAsEnvVars(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsEnvVars(&buf, vars))

	expected := `export TF_VAR_availability_zone_names='["us-west-1a"]'
export TF_VAR_aws_amis='{ eu-west-1 = "ami-b1cf19c6", us-east-1 = "ami-de7ab6b6", us-west-1 = "ami-3f75767a", us-west-2 = "ami-21f78e11" }'
export TF_VAR_docker_ports='[{ external = 8300, internal = 8301, protocol = "tcp" }]'
export TF_VAR_instance_name='my-instance'
export TF_VAR_password=''
export TF_VAR_region=''
`
	assert.Equal(t, expected, buf.String())
}

func TestWriteAsTFVars(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsTFVars(true, &buf, vars))

	expected := `availability_zone_names = ["us-west-1a"]
aws_amis = {
  eu-west-1 = "ami-b1cf19c6"
  us-east-1 = "ami-de7ab6b6"
  us-west-1 = "ami-3f75767a"
  us-west-2 = "ami-21f78e11"
}
docker_ports = [{
  external = 8300
  internal = 8301
  protocol = "tcp"
}]
instance_name = "my-instance"
# the root password to use with the database
password = null
region   = null
`
	assert.Equal(t, expected, buf.String())
}

func TestWriteAsTFEResource(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsTFEResource(&buf, vars))

	g := goldie.New(
		t,
		goldie.WithNameSuffix(".golden.tf"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	g.Assert(t, "tfe_resource", buf.Bytes())
}

func TestWriteAsWorkspacePayload(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsWorkspacePayload(&buf, vars))

	g := goldie.New(
		t,
		goldie.WithNameSuffix(".golden.json"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)

	g.Assert(t, "workspace_payload", buf.Bytes())
}

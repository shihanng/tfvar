package tfvar

import (
	"bytes"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/configs"
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
	assert.NoError(t, WriteAsTFVars(&buf, vars))

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
password      = null
region        = null
`
	assert.Equal(t, expected, buf.String())
}

func TestWriteAsTFEResource(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsTFEResource(&buf, vars))

	expected := `
resource "tfe_variable" "availability_zone_names" {
  key          = "availability_zone_names"
  value        = ["us-west-1a"]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "aws_amis" {
  key = "aws_amis"
  value = {
    eu-west-1 = "ami-b1cf19c6"
    us-east-1 = "ami-de7ab6b6"
    us-west-1 = "ami-3f75767a"
    us-west-2 = "ami-21f78e11"
  }
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "docker_ports" {
  key = "docker_ports"
  value = [{
    external = 8300
    internal = 8301
    protocol = "tcp"
  }]
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "instance_name" {
  key          = "instance_name"
  value        = "my-instance"
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "password" {
  key          = "password"
  value        = null
  sensitive    = true
  description  = "the root password to use with the database"
  workspace_id = null
  category     = "terraform"
}

resource "tfe_variable" "region" {
  key          = "region"
  value        = null
  sensitive    = false
  description  = ""
  workspace_id = null
  category     = "terraform"
}
`

	assert.Equal(t, expected, buf.String())
}

func TestWriteAsWorkspacePayload(t *testing.T) {
	vars, err := Load("testdata/defaults")
	require.NoError(t, err)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	var buf bytes.Buffer
	assert.NoError(t, WriteAsWorkspacePayload(&buf, vars))

	expected := `{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "availability_zone_names",
					"value":       "['us-west-1a']",
					"description": "",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   false
				}
			}
		}
		{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "aws_amis",
					"value":       "{ eu-west-1 = 'ami-b1cf19c6', us-east-1 = 'ami-de7ab6b6', us-west-1 = 'ami-3f75767a', us-west-2 = 'ami-21f78e11' }",
					"description": "",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   false
				}
			}
		}
		{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "docker_ports",
					"value":       "[{ external = 8300, internal = 8301, protocol = 'tcp' }]",
					"description": "",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   false
				}
			}
		}
		{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "instance_name",
					"value":       "my-instance",
					"description": "",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   false
				}
			}
		}
		{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "password",
					"value":       "",
					"description": "the root password to use with the database",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   true
				}
			}
		}
		{
			"data": {
				"type": "vars",
				"attributes": {
					"key":         "region",
					"value":       "",
					"description": "",
					"category":    "terraform",
					"hcl":         false,
					"sensitive":   false
				}
			}
		}
		`

	assert.Equal(t, expected, buf.String())
}

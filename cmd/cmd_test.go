package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlain(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = ["us-west-1a"]
docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
image_id                = null
`, actual.String())
}

func TestEnvVar(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata -e")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `export TF_VAR_availability_zone_names='["us-west-1a"]'
export TF_VAR_docker_ports='[{ external = 8300, internal = 8300, protocol = "tcp" }]'
export TF_VAR_image_id=''
`, actual.String())
}

func TestIgnoreDefault(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata --ignore-default")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = null
docker_ports            = null
image_id                = null
`, actual.String())
}

func TestAutoAssign(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata -a")
	os.Setenv("TF_VAR_image_id", "abc123")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = ["my-zone"]
docker_ports            = [{ external = 80, internal = 80, protocol = "tcp" }]
image_id                = "abc123"
`, actual.String())
}

func TestVar(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata -a --var='image_id=abc123' --var='unknown=xxx'")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = ["my-zone"]
docker_ports            = [{ external = 80, internal = 80, protocol = "tcp" }]
image_id                = "abc123"
`, actual.String())
}

func TestVarError(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata -a --var='unknown'")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	assert.Error(t, cmd.Execute())
	assert.Contains(t, actual.String(), `Error: tfvar: bad var string ''unknown''`)
}

func TestVarFile(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata --var-file testdata/my.tfvars")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = ["us-west-1a"]
docker_ports            = [{ external = 8300, internal = 8300, protocol = "tcp" }]
image_id                = "xyz"
`, actual.String())
}

func TestVarFileError(t *testing.T) {
	os.Args = strings.Fields("tfvar testdata --var-file testdata/bad.tfvars")

	var actual bytes.Buffer
	cmd, sync := New(&actual, "dev")
	defer sync()

	assert.Error(t, cmd.Execute())
	assert.Contains(t, actual.String(), `Error: tfvar: failed to parse 'testdata/bad.tfvars'`)
}


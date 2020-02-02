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
	cmd, sync := New(&actual)
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
	cmd, sync := New(&actual)
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
	cmd, sync := New(&actual)
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
	cmd, sync := New(&actual)
	defer sync()

	require.NoError(t, cmd.Execute())
	assert.Equal(t, `availability_zone_names = ["my-zone"]
docker_ports            = [{ external = 80, internal = 80, protocol = "tcp" }]
image_id                = "abc123"
`, actual.String())
}

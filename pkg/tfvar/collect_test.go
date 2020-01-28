package tfvar

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

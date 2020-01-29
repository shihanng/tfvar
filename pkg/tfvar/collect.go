package tfvar

import (
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type UnparsedVariableValue interface {
	ParseVariableValue(configs.VariableParsingMode) (cty.Value, error)
}

func collectFromEnvVars(to map[string]UnparsedVariableValue) {
	env := os.Environ()
	for _, raw := range env {
		if !strings.HasPrefix(raw, VarEnvPrefix) {
			continue
		}
		raw = raw[len(VarEnvPrefix):] // trim the prefix

		eq := strings.Index(raw, "=")

		// Igoring the one with no "="
		if eq > 0 {
			name := raw[:eq]
			rawVal := raw[eq+1:]

			to[name] = unparsedVariableValueString{
				str:  rawVal,
				name: name,
			}
		}
	}
}

func collectFromString(raw string, to map[string]UnparsedVariableValue) error {
	eq := strings.Index(raw, "=")
	if eq == -1 {
		return errors.Errorf("tfvar: bad var string '%s'", raw)
	}

	name := raw[:eq]
	rawVal := raw[eq+1:]

	to[name] = unparsedVariableValueString{
		str:  rawVal,
		name: name,
	}

	return nil
}

type unparsedVariableValueString struct {
	str  string
	name string
}

func (v unparsedVariableValueString) ParseVariableValue(mode configs.VariableParsingMode) (cty.Value, error) {
	val, hclDiags := mode.Parse(v.name, v.str)
	if hclDiags.HasErrors() {
		return cty.Value{}, errors.Wrap(hclDiags, "tfvar: failed to parse unparsedVariableValueString")
	}

	return val, nil
}

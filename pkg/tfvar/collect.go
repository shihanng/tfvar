package tfvar

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/json"
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

func collectFromFile(filename string, to map[string]UnparsedVariableValue) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Errorf("tfvar: reading file '%s'", filename)
	}

	var f *hcl.File
	if strings.HasSuffix(filename, ".json") {
		var hclDiags hcl.Diagnostics
		f, hclDiags = json.Parse(src, filename)
		if hclDiags.HasErrors() {
			return errors.Wrapf(hclDiags, "tfvar: failed to parse '%s'", filename)
		}
	} else {
		var hclDiags hcl.Diagnostics
		f, hclDiags = hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
		if hclDiags.HasErrors() {
			return errors.Wrapf(hclDiags, "tfvar: failed to parse '%s'", filename)
		}
	}

	attrs, hclDiags := f.Body.JustAttributes()
	if hclDiags.HasErrors() {
		return errors.Wrap(hclDiags, "tfvar: failed to get attributes")
	}

	for name, attr := range attrs {
		to[name] = unparsedVariableValueExpression{
			expr: attr.Expr,
		}
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

type unparsedVariableValueExpression struct {
	expr hcl.Expression
}

func (v unparsedVariableValueExpression) ParseVariableValue(_ configs.VariableParsingMode) (cty.Value, error) {
	val, hclDiags := v.expr.Value(nil) // nil because no function calls or variable references are allowed here
	if hclDiags.HasErrors() {
		return cty.Value{}, errors.Wrap(hclDiags, "tfvar: failed to parse unparsedVariableValueExpression")
	}

	return val, nil
}

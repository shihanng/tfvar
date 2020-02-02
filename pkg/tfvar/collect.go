package tfvar

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/json"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultVarsFilename     = "terraform.tfvars"
	defaultVarsFilenameJSON = defaultVarsFilename + `.json`
)

// LookupTFVarsFiles search for terraform.tfvars, terraform.tfvars.json, *.auto.tfvars, and *.auto.tfvars.json in dir. The value of dir is include in the returned value.
func LookupTFVarsFiles(dir string) []string {
	var files []string

	d := filepath.Join(dir, defaultVarsFilename)
	if _, err := os.Stat(d); err == nil {
		files = append(files, d)
	}

	dj := filepath.Join(dir, defaultVarsFilenameJSON)
	if _, err := os.Stat(dj); err == nil {
		files = append(files, dj)
	}

	if infos, err := ioutil.ReadDir(dir); err == nil {
		// "infos" is already sorted by name, so we just need to filter it here.
		for _, info := range infos {
			name := info.Name()
			if !isAutoVarFile(name) {
				continue
			}

			files = append(files, filepath.Join(dir, name))
		}
	}

	return files
}

// isAutoVarFile determines if the file ends with .auto.tfvars or .auto.tfvars.json
// https://github.com/hashicorp/terraform/blob/e9d0822b2a60f15653da0120607e74df1e116422/command/meta.go#L635-L638
func isAutoVarFile(path string) bool {
	return strings.HasSuffix(path, ".auto.tfvars") ||
		strings.HasSuffix(path, ".auto.tfvars.json")
}

// UnparsedVariableValue describes the value of variable definitions defined in tfvars files, environment variables, and raw string.
type UnparsedVariableValue interface {
	ParseVariableValue(configs.VariableParsingMode) (cty.Value, error)
}

// CollectFromEnvVars extracts the variable definitions from all environment variables that prefixed with TF_VAR_.
func CollectFromEnvVars(to map[string]UnparsedVariableValue) {
	env := os.Environ()
	for _, raw := range env {
		if !strings.HasPrefix(raw, varEnvPrefix) {
			continue
		}
		raw = raw[len(varEnvPrefix):] // trim the prefix

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

// CollectFromString extracts the variable definition from the given string.
func CollectFromString(raw string, to map[string]UnparsedVariableValue) error {
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

// CollectFromFile extracts the variable definitions from the given file.
func CollectFromFile(filename string, to map[string]UnparsedVariableValue) error {
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

// ParseValues assigns defined variables into the matching declared variables.
func ParseValues(from map[string]UnparsedVariableValue, vars []Variable) ([]Variable, error) {
	for i, v := range vars {
		unparsed, found := from[v.Name]
		if !found {
			continue
		}

		val, err := unparsed.ParseVariableValue(v.parsingMode)
		if err != nil {
			return nil, err
		}

		vars[i].Value = val
	}

	return vars, nil
}

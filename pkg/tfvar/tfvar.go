package tfvar

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name  string
	Value cty.Value

	parsingMode configs.VariableParsingMode
}

func Load(rootDir string) ([]Variable, error) {
	parser := configs.NewParser(nil)

	modules, diag := parser.LoadConfigDir(rootDir)
	if diag.HasErrors() {
		return nil, errors.Wrap(diag, "tfvar: loading config")
	}

	variables := make([]Variable, 0, len(modules.Variables))

	for _, v := range modules.Variables {
		variables = append(variables, Variable{
			Name:  v.Name,
			Value: v.Default,

			parsingMode: v.ParsingMode,
		})
	}

	return variables, nil
}

const varEnvPrefix = "TF_VAR_"

func WriteAsEnvVars(w io.Writer, vars []Variable) error {
	var we error

	for _, v := range vars {
		val := convertNull(v.Value)

		t := hclwrite.TokensForValue(val)
		b := t.Bytes()
		b = bytes.TrimPrefix(b, []byte(`"`))
		b = bytes.TrimSuffix(b, []byte(`"`))

		if we == nil {
			_, err := fmt.Fprintf(w, "export %s%s='%s'\n", varEnvPrefix, v.Name, string(b))
			we = errors.Wrap(err, "tfvar: unexpected writing export")
		}
	}

	return we
}

func WriteAsTFVars(w io.Writer, vars []Variable) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, v := range vars {
		rootBody.SetAttributeValue(v.Name, v.Value)
	}

	_, err := f.WriteTo(w)
	return errors.Wrap(err, "tfvar: failed to write as tfvars")
}

func convertNull(v cty.Value) cty.Value {
	if v.IsNull() {
		return cty.StringVal("")
	}

	return v
}

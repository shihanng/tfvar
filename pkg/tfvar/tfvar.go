// Package tfvar contains the essential tools to extract input variables from Terraform configurations, retrieve variable definitions from sources, and parse those values back into the input variables.
package tfvar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/shihanng/tfvar/pkg/configs"
	"github.com/zclconf/go-cty/cty"
)

// Variable represents a simplified version of Terraform's input variable, e.g.
//    variable "image_id" {
//      type = string
//    }
type Variable struct {
	Name        string
	Value       cty.Value
	Description string
	Sensitive   bool

	parsingMode configs.VariableParsingMode
}

// Load extracts all input variables declared in the Terraform configurations located in dir.
func Load(dir string) ([]Variable, error) {
	parser := configs.NewParser(nil)

	modules, diag := parser.LoadConfigDir(dir)
	if diag.HasErrors() {
		return nil, errors.Wrap(diag, "tfvar: loading config")
	}

	variables := make([]Variable, 0, len(modules.Variables))

	for _, v := range modules.Variables {
		variables = append(variables, Variable{
			Name:        v.Name,
			Value:       v.Default,
			Description: v.Description,
			Sensitive:   v.Sensitive,

			parsingMode: v.ParsingMode,
		})
	}

	return variables, nil
}

const varEnvPrefix = "TF_VAR_"

// WriteAsEnvVars outputs the given vars in environment variables format, e.g.
//    export TF_VAR_region='ap-northeast-1'
func WriteAsEnvVars(w io.Writer, vars []Variable) error {
	for _, v := range vars {
		val := convertNull(v.Value)

		t := hclwrite.TokensForValue(val)
		t = oneliner(t)
		b := hclwrite.Format(t.Bytes())
		b = bytes.TrimPrefix(b, []byte(`"`))
		b = bytes.TrimSuffix(b, []byte(`"`))

		if _, err := fmt.Fprintf(w, "export %s%s='%s'\n", varEnvPrefix, v.Name, string(b)); err != nil {
			return errors.Wrap(err, "tfvar: unexpected writing export")
		}
	}

	return nil
}

func oneliner(original hclwrite.Tokens) hclwrite.Tokens {
	var toks hclwrite.Tokens

	for i, t := range original {
		if t.Type != hclsyntax.TokenNewline {
			toks = append(toks, t)
			continue
		}

		// https://github.com/hashicorp/hcl/blob/v2.6.0/hclwrite/generate.go#L117-L156
		// Newline only exists in map/object type (between hclsyntax.TokenOBrace and hclsyntax.TokenCBrace).
		if original[i-1].Type == hclsyntax.TokenOBrace || original[i+1].Type == hclsyntax.TokenCBrace {
			continue
		}

		// Replace newline with comma.
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenComma,
			Bytes: []byte{','},
		})
	}

	return toks
}

// WriteAsTFVars outputs the given vars in Terraform's variable definitions format, e.g.
//    region = "ap-northeast-1"
func WriteAsTFVars(withComments bool, w io.Writer, vars []Variable) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, v := range vars {
		if withComments {
			comment := strings.TrimSpace(v.Description)
			if len(comment) > 0 {
				rootBody.AppendUnstructuredTokens(commentToTokens(comment))
			}
		}
		rootBody.SetAttributeValue(v.Name, v.Value)
	}

	_, err := f.WriteTo(w)
	return errors.Wrap(err, "tfvar: failed to write as tfvars")
}

func commentToTokens(comment string) hclwrite.Tokens {
	lines := strings.Split(comment, "\n")
	newLine := hclwrite.Token{
		Type:  hclsyntax.TokenNewline,
		Bytes: []byte(string(hclsyntax.TokenNewline)),
	}

	var tokens hclwrite.Tokens
	for _, line := range lines {
		tokens = append(tokens, &hclwrite.Token{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(fmt.Sprintf("# %s", line)),
		})
		tokens = append(tokens, &newLine)
	}

	return tokens
}

type workspacePayload struct {
	Data workspaceData `json:"data"`
}

type workspaceData struct {
	Type       string              `json:"type"`
	Attributes workspaceAttributes `json:"attributes"`
}

type workspaceAttributes struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Category    string `json:"category"`
	HCL         bool   `json:"hcl"`
	Sensitive   bool   `json:"sensitive"`
}

func WriteAsWorkspacePayload(w io.Writer, vars []Variable) error {
	for _, v := range vars {
		val := convertNull(v.Value)

		t := hclwrite.TokensForValue(val)
		t = oneliner(t)
		b := hclwrite.Format(t.Bytes())
		b = bytes.TrimPrefix(b, []byte(`"`))
		b = bytes.TrimSuffix(b, []byte(`"`))
		b = bytes.ReplaceAll(b, []byte(`"`), []byte(`'`))

		payload := workspacePayload{
			Data: workspaceData{
				Type: "vars",
				Attributes: workspaceAttributes{
					Key:         v.Name,
					Value:       string(b),
					Description: v.Description,
					Category:    "terraform",
					HCL:         false,
					Sensitive:   v.Sensitive,
				},
			},
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(payload); err != nil {
			return errors.Wrap(err, "tfvar: unexpected error writing payload")
		}
	}

	return nil
}

func WriteAsTFEResource(w io.Writer, vars []Variable) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, v := range vars {
		rootBody.AppendNewline()
		resourceBlock := rootBody.AppendNewBlock("resource", []string{"tfe_variable", v.Name})
		resourceBody := resourceBlock.Body()
		resourceBody.SetAttributeValue("key", cty.StringVal(v.Name))
		resourceBody.SetAttributeValue("value", v.Value)
		resourceBody.SetAttributeValue("sensitive", cty.BoolVal(v.Sensitive))
		resourceBody.SetAttributeValue("description", cty.StringVal(v.Description))
		resourceBody.SetAttributeValue("workspace_id", cty.NilVal)
		resourceBody.SetAttributeValue("category", cty.StringVal("terraform"))
	}

	_, err := f.WriteTo(w)
	return errors.Wrap(err, "tfe_variable: failed to write as tfe_variable resource")
}

func convertNull(v cty.Value) cty.Value {
	if v.IsNull() {
		return cty.StringVal("")
	}

	return v
}

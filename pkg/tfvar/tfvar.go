package tfvar

import (
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name  string
	Value cty.Value
}

func Load(rootDir string) ([]Variable, error) {
	parser := configs.NewParser(nil)

	modules, err := parser.LoadConfigDir(rootDir)
	if err != nil {
		return nil, errors.Wrap(err, "tfvar: loading config")
	}

	variables := make([]Variable, 0, len(modules.Variables))

	for _, v := range modules.Variables {
		variables = append(variables, Variable{
			Name:  v.Name,
			Value: v.Default,
		})
	}

	return variables, nil
}

package tfvar

import (
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform/configs"
)

func Load(rootDir string) ([]string, error) {
	parser := configs.NewParser(nil)

	modules, err := parser.LoadConfigDir(rootDir)
	if err != nil {
		return nil, errors.Wrap(err, "tfvar: loading config")
	}

	names := make([]string, 0, len(modules.Variables))

	for _, v := range modules.Variables {
		names = append(names, v.Name)
	}

	return names, nil
}

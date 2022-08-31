package configs

import "github.com/hashicorp/hcl/v2"

var moduleBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "source",
			Required: true,
		},
		{
			Name: "version",
		},
		{
			Name: "count",
		},
		{
			Name: "for_each",
		},
		{
			Name: "depends_on",
		},
		{
			Name: "providers",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "_"}, // meta-argument escaping block

		// These are all reserved for future use.
		{Type: "lifecycle"},
		{Type: "locals"},
		{Type: "provider", LabelNames: []string{"type"}},
	},
}

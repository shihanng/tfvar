// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// Module is a container for a set of configuration constructs that are
// evaluated within a common namespace.
type Module struct {
	Variables map[string]*Variable
}

// File describes the contents of a single configuration file.
//
// Individual files are not usually used alone, but rather combined together
// with other files (conventionally, those in the same directory) to produce
// a *Module, using NewModule.
//
// At the level of an individual file we represent directly the structural
// elements present in the file, without any attempt to detect conflicting
// declarations. A File object can therefore be used for some basic static
// analysis of individual elements, but must be built into a Module to detect
// duplicate declarations.
type File struct {
	Variables []*Variable
}

// NewModule takes a list of primary files and a list of override files and
// produces a *Module by combining the files together.
//
// If there are any conflicting declarations in the given files -- for example,
// if the same variable name is defined twice -- then the resulting module
// will be incomplete and error diagnostics will be returned. Careful static
// analysis of the returned Module is still possible in this case, but the
// module will probably not be semantically valid.
func NewModule(primaryFiles, overrideFiles []*File) (*Module, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	mod := &Module{
		Variables: map[string]*Variable{},
	}

	for _, file := range primaryFiles {
		fileDiags := mod.appendFile(file)
		diags = append(diags, fileDiags...)
	}

	for _, file := range overrideFiles {
		fileDiags := mod.mergeFile(file)
		diags = append(diags, fileDiags...)
	}

	return mod, diags
}

func (m *Module) appendFile(file *File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range file.Variables {
		if existing, exists := m.Variables[v.Name]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable declaration",
				Detail:   fmt.Sprintf("A variable named %q was already declared at %s. Variable names must be unique within a module.", existing.Name, existing.DeclRange),
				Subject:  &v.DeclRange,
			})
		}
		m.Variables[v.Name] = v
	}

	return diags
}

func (m *Module) mergeFile(file *File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range file.Variables {
		existing, exists := m.Variables[v.Name]
		if !exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing base variable declaration to override",
				Detail:   fmt.Sprintf("There is no variable named %q. An override file can only override a variable that was already declared in a primary configuration file.", v.Name),
				Subject:  &v.DeclRange,
			})
			continue
		}
		mergeDiags := existing.merge(v)
		diags = append(diags, mergeDiags...)
	}

	return diags
}

package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func (v *Variable) merge(ov *Variable) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if ov.DescriptionSet {
		v.Description = ov.Description
		v.DescriptionSet = ov.DescriptionSet
	}
	if ov.SensitiveSet {
		v.Sensitive = ov.Sensitive
		v.SensitiveSet = ov.SensitiveSet
	}
	if ov.Default != cty.NilVal {
		v.Default = ov.Default
	}
	if ov.Type != cty.NilType {
		v.Type = ov.Type
		v.ConstraintType = ov.ConstraintType
	}
	if ov.ParsingMode != 0 {
		v.ParsingMode = ov.ParsingMode
	}
	if ov.NullableSet {
		v.Nullable = ov.Nullable
		v.NullableSet = ov.NullableSet
	}

	// If the override file overrode type without default or vice-versa then
	// it may have created an invalid situation, which we'll catch now by
	// attempting to re-convert the value.
	//
	// Note that here we may be re-converting an already-converted base value
	// from the base config. This will be a no-op if the type was not changed,
	// but in particular might be user-observable in the edge case where the
	// literal value in config could've been converted to the overridden type
	// constraint but the converted value cannot. In practice, this situation
	// should be rare since most of our conversions are interchangable.
	if v.Default != cty.NilVal {
		val, err := convert.Convert(v.Default, v.ConstraintType)
		if err != nil {
			// What exactly we'll say in the error message here depends on whether
			// it was Default or Type that was overridden here.
			switch {
			case ov.Type != cty.NilType && ov.Default == cty.NilVal:
				// If only the type was overridden
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("Overriding this variable's type constraint has made its default value invalid: %s.", err),
					Subject:  &ov.DeclRange,
				})
			case ov.Type == cty.NilType && ov.Default != cty.NilVal:
				// Only the default was overridden
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("The overridden default value for this variable is not compatible with the variable's type constraint: %s.", err),
					Subject:  &ov.DeclRange,
				})
			default:
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid default value for variable",
					Detail:   fmt.Sprintf("This variable's default value is not compatible with its type constraint: %s.", err),
					Subject:  &ov.DeclRange,
				})
			}
		} else {
			v.Default = val
		}

		// ensure a null default wasn't merged in when it is not allowed
		if !v.Nullable && v.Default.IsNull() {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid default value for variable",
				Detail:   "A null default value is not valid when nullable=false.",
				Subject:  &ov.DeclRange,
			})
		}
	}

	return diags
}

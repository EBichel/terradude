package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Global struct {
	Name  string
	Value cty.Value
}

func DecodeGlobalCty(hclconfigs []*Config, ctx *hcl.EvalContext) (*cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	globals := map[string]cty.Value{}

	for _, hclconfig := range hclconfigs {
		if hclconfig.Globals != nil {
			attrs, diag := hclconfig.Globals.Body.JustAttributes()
			if diag.HasErrors() {
				return nil, diags
			}
			diags = append(diags, diag...)
			for _, attr := range attrs {
				if _, ok := globals[attr.Name]; ok {
					continue
				}
				value, diag := attr.Expr.Value(ctx)
				diags = append(diags, diag...)
				if diag.HasErrors() {
					return nil, diags
				}
				globals[attr.Name] = value
			}
		}
	}

	ctyGlobals, err := mapToCty(globals)
	if err != nil {
		return nil, err.(hcl.Diagnostics)
	}

	return ctyGlobals, nil
}

func mapToCty(theMap map[string]cty.Value) (*cty.Value, error) {
	ctyTypes := map[string]cty.Type{}
	for key, value := range theMap {
		ctyTypes[key] = value.Type()
	}
	ctyObject := cty.Object(ctyTypes)
	ctyGlobals, err := gocty.ToCtyValue(theMap, ctyObject)
	if err != nil {
		return nil, err
	}
	return &ctyGlobals, nil
}

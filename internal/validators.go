package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ValidateSetInSlice struct {
	values []string
}

func (r ValidateSetInSlice) Description(ctx context.Context) string {
	return ""
}

func (r ValidateSetInSlice) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}

func (r ValidateSetInSlice) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {

	var events types.Set
	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &events)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if events.IsNull() || events.IsUnknown() {
		return
	}

	var elems []types.String
	diags = events.ElementsAs(ctx, &elems, false)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, attr_value := range elems {
		r.validateAttr(attr_value, req, resp)

		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (r ValidateSetInSlice) validateAttr(attr types.String, req validator.SetRequest, resp *validator.SetResponse) {

	str := attr.ValueString()

	for _, value := range r.values {
		if value == str {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Element",
		fmt.Sprintf("%s is not supported (expected: %s)", attr, strings.Join(r.values, ", ")),
	)
}

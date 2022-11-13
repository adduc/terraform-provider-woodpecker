package internal

import (
	"context"
	"fmt"
	"strings"

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

func (r ValidateSetInSlice) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

	var events types.Set
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &events)

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

func (r ValidateSetInSlice) validateAttr(attr types.String, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

	str := attr.ValueString()

	for _, value := range r.values {
		if value == str {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid Element",
		fmt.Sprintf("%s is not supported (expected: %s)", attr, strings.Join(r.values, ", ")),
	)
}

package internal

import "github.com/hashicorp/terraform-plugin-framework/types"

type Repo struct {
	ID         types.Int64  `tfsdk:"id"`
	Owner      types.String `tfsdk:"owner"`
	Name       types.String `tfsdk:"name"`
	FullName   types.String `tfsdk:"full_name"`
	Avatar     types.String `tfsdk:"avatar"`
	Link       types.String `tfsdk:"link"`
	Kind       types.String `tfsdk:"kind"`
	Clone      types.String `tfsdk:"clone"`
	Branch     types.String `tfsdk:"branch"`
	Timeout    types.Int64  `tfsdk:"timeout"`
	Visibility types.String `tfsdk:"visibility"`
	IsPrivate  types.Bool   `tfsdk:"is_private"`
	IsTrusted  types.Bool   `tfsdk:"is_trusted"`
	IsGated    types.Bool   `tfsdk:"is_gated"`
	AllowPull  types.Bool   `tfsdk:"allow_pull"`
	Config     types.String `tfsdk:"config"`
}

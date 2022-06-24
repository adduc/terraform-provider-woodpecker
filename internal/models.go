package internal

import "github.com/hashicorp/terraform-plugin-framework/types"

type Repo struct {
	ID    types.Int64  `tfsdk:"id"`
	Owner types.String `tfsdk:"owner"`
	Name  types.String `tfsdk:"name"`
}

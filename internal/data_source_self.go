package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDataSourceSelf() datasource.DataSource {
	return &DataSourceSelf{}
}

type DataSourceSelf struct {
	p woodpeckerProvider
}

func (d *DataSourceSelf) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_self"
}

func (r DataSourceSelf) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "User ID",
			},
			"login": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Username for user",
			},
			"email": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Email address for user",
			},
			"avatar": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Avatar URL for user",
			},
			"active": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether user is active in the system",
			},
			"admin": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether user is a Woodpecker admin",
			},
		},
	}, nil
}

func (r *DataSourceSelf) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*woodpeckerProvider)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *woodpeckerProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.p = *p
}

func (r DataSourceSelf) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var resourceData User
	self := r.p.self

	// unmarshall self response into resourceData
	resourceData.ID = types.Int64{Value: self.ID}
	resourceData.Login = types.String{Value: self.Login}
	resourceData.Email = types.String{Value: self.Email}
	resourceData.Avatar = types.String{Value: self.Avatar}
	resourceData.Active = types.Bool{Value: self.Active}
	resourceData.Admin = types.Bool{Value: self.Admin}

	diags := resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDataSourceUser() datasource.DataSource {
	return &DataSourceUser{}
}

type DataSourceUser struct {
	p woodpeckerProvider
}

func (d *DataSourceUser) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r DataSourceUser) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Use this data source to get information on an existing user",

		Attributes: map[string]tfsdk.Attribute{

			// Required Attributes
			"login": {
				Type:        types.StringType,
				Required:    true,
				Description: "Username for user",
			},

			// Computed Attributes
			"email": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Email address for user",
			},
			"active": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether user is active in the system",
			},
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "User ID",
			},
			"admin": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "Whether user is a Woodpecker admin",
			},
			"avatar": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Avatar URL for user",
			},
		},
	}, nil
}

func (r *DataSourceUser) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r DataSourceUser) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// unmarshall request config into resourceData
	var resourceData User
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	login := resourceData.Login.ValueString()

	user, err := r.p.client.User(login)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	WoodpeckerToUser(ctx, *user, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

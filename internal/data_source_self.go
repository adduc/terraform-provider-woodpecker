package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceSelfType struct{}

func (r dataSourceSelfType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r dataSourceSelfType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceSelf{
		p: *(p.(*woodpeckerProvider)),
	}, nil
}

type dataSourceSelf struct {
	p woodpeckerProvider
}

func (r dataSourceSelf) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

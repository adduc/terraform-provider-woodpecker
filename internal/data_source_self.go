package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (r DataSourceSelf) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on the authenticated user used by Terraform",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "User ID",
			},
			"login": schema.StringAttribute{
				Computed:    true,
				Description: "Username for user",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Email address for user",
			},
			"avatar": schema.StringAttribute{
				Computed:    true,
				Description: "Avatar URL for user",
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether user is active in the system",
			},
			"admin": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether user is a Woodpecker admin",
			},
		},
	}
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
	resourceData.ID = types.Int64Value(self.ID)
	resourceData.Login = types.StringValue(self.Login)
	resourceData.Email = types.StringValue(self.Email)
	resourceData.Avatar = types.StringValue(self.Avatar)
	resourceData.Active = types.BoolValue(self.Active)
	resourceData.Admin = types.BoolValue(self.Admin)

	diags := resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

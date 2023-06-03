package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewDataSourceSecret() datasource.DataSource {
	return &DataSourceSecret{}
}

type DataSourceSecret struct {
	client woodpecker.Client
}

func (d *DataSourceSecret) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r DataSourceSecret) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on an existing global secret",

		Attributes: map[string]schema.Attribute{
			// Required Attributes
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Secret Name",
			},

			// Computed Attributes
			"plugins_only": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether secret is only available for plugins",
			},
			"images": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of images where this secret is available, leave empty to allow all images",
			},
			"events": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "One or more event types where secret is available (push, tag, pull_request, deployment, cron, manual)",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
		},
	}
}

func (r *DataSourceSecret) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = p.client
}

func (r DataSourceSecret) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData SecretData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	secretName := resourceData.Name.ValueString()

	secret, err := r.client.GlobalSecret(secretName)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving secret", err.Error())
		return
	}

	diags = r.WoodpeckerToSecretData(ctx, *secret, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r DataSourceSecret) WoodpeckerToSecretData(ctx context.Context, wSecret woodpecker.Secret, secret *SecretData) diag.Diagnostics {

	var diags, err diag.Diagnostics

	secret.ID = types.Int64Value(wSecret.ID)
	secret.Name = types.StringValue(wSecret.Name)
	secret.PluginsOnly = types.BoolValue(wSecret.PluginsOnly)
	secret.Images, diags = types.SetValueFrom(ctx, types.StringType, wSecret.Images)
	secret.Events, err = types.SetValueFrom(ctx, types.StringType, wSecret.Events)

	diags.Append(err...)

	return diags
}

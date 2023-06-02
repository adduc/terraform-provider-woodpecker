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

func NewDataSourceRepositoryRegistry() datasource.DataSource {
	return &DataSourceRepositoryRegistry{}
}

type DataSourceRepositoryRegistry struct {
	client woodpecker.Client
}

func (d *DataSourceRepositoryRegistry) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_registry"
}

func (r DataSourceRepositoryRegistry) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on an existing registry for a repository",

		Attributes: map[string]schema.Attribute{
			// Required Attributes
			"repo_owner": schema.StringAttribute{
				Required:    true,
				Description: "User or organization responsible for repository",
			},
			"repo_name": schema.StringAttribute{
				Required:    true,
				Description: "Repository name",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Registry Address",
			},

			// Computed Attributes
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "Registry Username",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Description: "Registry Token",
				Sensitive:   true,
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Registry Email",
			},
		},
	}
}

func (r *DataSourceRepositoryRegistry) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r DataSourceRepositoryRegistry) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryRegistryData
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()
	address := resourceData.Address.ValueString()

	registry, err := r.client.Registry(repoOwner, repoName, address)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving repository secret", err.Error())
		return
	}

	diags = r.WoodpeckerToRepositoryRegistryData(ctx, *registry, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r DataSourceRepositoryRegistry) WoodpeckerToRepositoryRegistryData(ctx context.Context, wRegistry woodpecker.Registry, registry *RepositoryRegistryData) diag.Diagnostics {

	var diags diag.Diagnostics

	registry.ID = types.Int64Value(wRegistry.ID)
	registry.Address = types.StringValue(wRegistry.Address)
	registry.Username = types.StringValue(wRegistry.Username)
	registry.Token = types.StringValue(wRegistry.Token)
	registry.Email = types.StringValue(wRegistry.Email)

	return diags
}

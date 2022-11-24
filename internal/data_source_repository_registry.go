package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func (r DataSourceRepositoryRegistry) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Use this data source to get information on an existing registry for a repository",

		Attributes: map[string]tfsdk.Attribute{
			// Required Attributes
			"repo_owner": {
				Type:        types.StringType,
				Required:    true,
				Description: "User or organization responsible for repository",
			},
			"repo_name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Repository name",
			},
			"address": {
				Type:        types.StringType,
				Required:    true,
				Description: "Registry Address",
			},

			// Computed Attributes
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"username": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Registry Username",
			},
			"token": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Registry Token",
				Sensitive:   true,
			},
			"email": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Registry Email",
			},
		},
	}, nil
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

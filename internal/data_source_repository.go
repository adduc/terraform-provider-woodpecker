package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDataSourceRepository() datasource.DataSource {
	return &DataSourceRepository{}
}

type DataSourceRepository struct {
	p woodpeckerProvider
}

func (d *DataSourceRepository) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r DataSourceRepository) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Use this data source to get information on an existing repository",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "Repository ID",
			},
			"owner": {
				Type:        types.StringType,
				Required:    true,
				Description: "User or organization responsible for repository",
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Repository name",
			},
			"full_name": {
				Type:        types.StringType,
				Computed:    true,
				Description: "*owner*/*name*",
			},
			"avatar": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Repository avatar URL",
			},
			"link": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Link to repository",
			},
			"kind": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Kind of repository (e.g. git)",
			},
			"clone": {
				Type:        types.StringType,
				Computed:    true,
				Description: "URL to clone repository",
			},
			"branch": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Default branch name",
			},
			"timeout": {
				Type:     types.Int64Type,
				Computed: true,
				Description: "After this timeout (in minutes) a pipeline has " +
					"to finish or will be treated as timed out.",
			},
			"visibility": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Public, Private, or Internal",
			},
			"is_trusted": {
				Type:     types.BoolType,
				Computed: true,
				Description: "If true, underlying pipeline containers get " +
					"access to escalated capabilities like mounting volumes.",
			},
			"is_gated": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "When true, every pipeline needs to be approved before being executed.",
			},
			"allow_pull": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "If true, pipelines can run on pull requests.",
			},
			"config": {
				Type:     types.StringType,
				Computed: true,
				MarkdownDescription: "Path to the pipeline config file or " +
					"folder. When empty, defaults to `.woodpecker/*.yml` -> " +
					"`.woodpecker.yml` -> `.drone.yml`.",
			},
		},
	}, nil
}

func (r *DataSourceRepository) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r DataSourceRepository) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// unmarshall request config into resourceData
	var resourceData Repository
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.Owner.ValueString()
	repoName := resourceData.Name.ValueString()

	repo, err := r.p.client.Repo(repoOwner, repoName)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving repo", err.Error())
		return
	}

	WoodpeckerToRepository(*repo, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

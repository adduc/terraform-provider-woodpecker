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

func NewDataSourceRepositoryCron() datasource.DataSource {
	return &DataSourceRepositoryCron{}
}

type DataSourceRepositoryCron struct {
	client woodpecker.Client
}

func (d *DataSourceRepositoryCron) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_cron"
}

func (r DataSourceRepositoryCron) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
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
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Cron Name",
			},

			// Computed Attributes
			"branch": {
				Type:        types.StringType,
				Computed:    true,
				Description: "",
			},
			"created": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"creator_id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"next_exec": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"repo_id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"schedule": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Schedule (based on UTC)",
			},
		},
	}, nil
}

func (r *DataSourceRepositoryCron) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r DataSourceRepositoryCron) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryCron
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()
	cronId := resourceData.ID.ValueInt64()

	cron, err := r.client.CronGet(repoOwner, repoName, cronId)

	if err != nil {
		resp.Diagnostics.AddError("Error retrieving cron", err.Error())
		return
	}

	WoodpeckerToRepositoryCron(*cron, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

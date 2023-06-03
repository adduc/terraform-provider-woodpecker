package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (r DataSourceRepositoryCron) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on an existing cron for a repository",

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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Cron Name",
			},

			// Computed Attributes
			"branch": schema.StringAttribute{
				Computed:    true,
				Description: "",
			},
			"created": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"creator_id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"next_exec": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"repo_id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"schedule": schema.StringAttribute{
				Computed:    true,
				Description: "Schedule (based on UTC)",
			},
		},
	}
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
		resp.Diagnostics.AddError("Error retrieving repository cron", err.Error())
		return
	}

	WoodpeckerToRepositoryCron(*cron, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceRepoType struct{}

func (r dataSourceRepoType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
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
			"is_private": {
				Type:     types.BoolType,
				Computed: true,
				Description: "If true, only authenticated users of the " +
					"Woodpecker instance can see this project.",
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

func (r dataSourceRepoType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceRepo{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceRepo struct {
	p provider
}

func (r dataSourceRepo) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {

	// unmarshall request config into resourceData
	var resourceData Repo
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.Owner.Value
	repoName := resourceData.Name.Value

	repo, err := r.p.client.Repo(repoOwner, repoName)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving repo", err.Error())
		return
	}

	// unmarshall repo response into resourceData
	resourceData.ID = types.Int64{Value: repo.ID}
	resourceData.FullName = types.String{Value: repo.FullName}
	resourceData.Avatar = types.String{Value: repo.Avatar}
	resourceData.Link = types.String{Value: repo.Link}
	resourceData.Kind = types.String{Value: repo.Kind}
	resourceData.Clone = types.String{Value: repo.Clone}
	resourceData.Branch = types.String{Value: repo.Branch}
	resourceData.Timeout = types.Int64{Value: repo.Timeout}
	resourceData.Visibility = types.String{Value: repo.Visibility}
	resourceData.IsPrivate = types.Bool{Value: repo.IsPrivate}
	resourceData.IsTrusted = types.Bool{Value: repo.IsTrusted}
	resourceData.IsGated = types.Bool{Value: repo.IsGated}
	resourceData.AllowPull = types.Bool{Value: repo.AllowPull}
	resourceData.Config = types.String{Value: repo.Config}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

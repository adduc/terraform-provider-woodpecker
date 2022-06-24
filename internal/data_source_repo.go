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
				Type:     types.Int64Type,
				Computed: true,
			},
			"owner": {
				Type:     types.StringType,
				Required: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
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
	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

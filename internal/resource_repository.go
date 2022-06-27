package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type resourceRepositoryType struct{}

func (r resourceRepositoryType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			// Required Attributes
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

			// Optional Attributes
			"timeout": {
				Type:     types.Int64Type,
				Optional: true,
				Computed: true,
				Description: "After this timeout (in minutes) a pipeline has " +
					"to finish or will be treated as timed out.",
			},
			"visibility": {
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Public, Private, or Internal",
			},
			"is_trusted": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				Description: "If true, underlying pipeline containers get " +
					"access to escalated capabilities like mounting volumes.",
			},
			"is_gated": {
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				Description: "When true, every pipeline needs to be approved before being executed.",
			},
			"allow_pull": {
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				Description: "If true, pipelines can run on pull requests.",
			},
			"config": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				MarkdownDescription: "Path to the pipeline config file or " +
					"folder. When empty, defaults to `.woodpecker/*.yml` -> " +
					"`.woodpecker.yml` -> `.drone.yml`.",
			},

			// Computed Attributes
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "Repository ID",
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
		},
	}, nil
}

func (r resourceRepositoryType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceRepository{
		p: *(p.(*provider)),
	}, nil
}

type resourceRepository struct {
	p provider
}

func (r resourceRepository) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// unmarshall request config into resourceData
	var resourceData Repository
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.Owner.Value
	repoName := resourceData.Name.Value

	_, err := r.p.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not fetch repository information", err.Error())
		return
	}

	_, err = r.p.client.RepoPost(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not activate repository", err.Error())
		return
	}

	patch := prepareRepositoryPatch(resourceData)

	_, err = r.p.client.RepoPatch(repoOwner, repoName, patch)

	if err != nil {
		resp.Diagnostics.AddError("Could not update reposiotry", err.Error())
		return
	}

	repo, err := r.p.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not refresh reposiotry", err.Error())
		return
	}

	WoodpeckerToRepository(*repo, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r resourceRepository) ModifyPlan(ctx context.Context, req tfsdk.ModifyResourcePlanRequest, resp *tfsdk.ModifyResourcePlanResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state Repository
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	plan.FullName = state.FullName
	plan.Avatar = state.Avatar

	plan.Link = state.Link
	plan.Kind = state.Kind
	plan.Clone = state.Clone
	plan.Branch = state.Branch

	if plan.Visibility.Unknown {
		plan.Visibility = state.Visibility
	}

	if plan.IsGated.Unknown {
		plan.IsGated = state.IsGated
	}

	if plan.IsTrusted.Unknown {
		plan.IsTrusted = state.IsTrusted
	}

	if plan.AllowPull.Unknown {
		plan.AllowPull = state.AllowPull
	}

	if plan.Timeout.Unknown {
		plan.Timeout = state.Timeout
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r resourceRepository) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// unmarshall request config into resourceData
	var resourceData Repository
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.Owner.Value
	repoName := resourceData.Name.Value

	repo, err := r.p.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToRepository(*repo, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r resourceRepository) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

	var repoPlan Repository
	diags := req.Plan.Get(ctx, &repoPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoState Repository
	diags = req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.Owner.Value
	repoName := repoState.Name.Value

	patch := prepareRepositoryPatch(repoPlan)

	repo, err := r.p.client.RepoPatch(repoOwner, repoName, patch)

	if err != nil {
		resp.Diagnostics.AddError("Could not update reposiotry", err.Error())
		return
	}

	WoodpeckerToRepository(*repo, &repoState)

	diags = resp.State.Set(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
}

func (r resourceRepository) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	var repoState Repository
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.Owner.Value
	repoName := repoState.Name.Value

	err := r.p.client.RepoDel(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r resourceRepository) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: owner/name. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("name"), idParts[1])...)
}

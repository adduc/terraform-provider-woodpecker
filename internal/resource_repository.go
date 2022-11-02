package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewRepositoryResource() resource.Resource {
	return &ResourceRepository{}
}

type ResourceRepository struct {
	client woodpecker.Client
}

func (r ResourceRepository) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r ResourceRepository) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *ResourceRepository) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceRepository) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	_, err := r.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not fetch repository information", err.Error())
		return
	}

	_, err = r.client.RepoPost(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not activate repository", err.Error())
		return
	}

	patch := prepareRepositoryPatch(resourceData)

	_, err = r.client.RepoPatch(repoOwner, repoName, patch)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository", err.Error())
		return
	}

	repo, err := r.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not refresh repository", err.Error())
		return
	}

	WoodpeckerToRepository(*repo, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepository) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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

	if plan.Visibility.IsUnknown() {
		plan.Visibility = state.Visibility
	}

	if plan.IsGated.IsUnknown() {
		plan.IsGated = state.IsGated
	}

	if plan.IsTrusted.IsUnknown() {
		plan.IsTrusted = state.IsTrusted
	}

	if plan.AllowPull.IsUnknown() {
		plan.AllowPull = state.AllowPull
	}

	if plan.Timeout.IsUnknown() {
		plan.Timeout = state.Timeout
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepository) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData Repository
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.Owner.ValueString()
	repoName := resourceData.Name.ValueString()

	repo, err := r.client.Repo(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToRepository(*repo, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepository) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

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

	repoOwner := repoState.Owner.ValueString()
	repoName := repoState.Name.ValueString()

	patch := prepareRepositoryPatch(repoPlan)

	repo, err := r.client.RepoPatch(repoOwner, repoName, patch)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository", err.Error())
		return
	}

	WoodpeckerToRepository(*repo, &repoState)

	diags = resp.State.Set(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepository) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var repoState Repository
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.Owner.ValueString()
	repoName := repoState.Name.ValueString()

	err := r.client.RepoDel(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceRepository) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: owner/name. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}

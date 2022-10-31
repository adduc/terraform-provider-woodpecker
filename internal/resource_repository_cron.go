package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewRepositoryCronResource() resource.Resource {
	return &ResourceRepositoryCron{}
}

type ResourceRepositoryCron struct {
	client woodpecker.Client
}

func (r ResourceRepositoryCron) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_cron"
}

func (r ResourceRepositoryCron) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"schedule": {
				Type:        types.StringType,
				Required:    true,
				Description: "Schedule (based on UTC)",
			},

			// Optional Attributes
			"repo_id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"creator_id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"next_exec": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"branch": {
				Type:        types.StringType,
				Computed:    true,
				Description: "",
			},

			// Computed
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
			"created": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
		},
	}, nil
}

func (r *ResourceRepositoryCron) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceRepositoryCron) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryCron
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()

	cron := prepareRepositoryCronPatch(resourceData)

	cron, err := r.client.CronCreate(repoOwner, repoName, cron)

	if err != nil {
		resp.Diagnostics.AddError("Could not create repository cron", err.Error())
		return
	}

	resourceData.RepoOwner = types.String{Value: repoOwner}
	resourceData.RepoName = types.String{Value: repoName}
	WoodpeckerToRepositoryCron(*cron, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryCron) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
}

func (r ResourceRepositoryCron) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryCron
	diags := req.State.Get(ctx, &resourceData)
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
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToRepositoryCron(*cron, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryCron) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var repoCronPlan RepositoryCron
	diags := req.Plan.Get(ctx, &repoCronPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoCronState RepositoryCron
	diags = req.State.Get(ctx, &repoCronState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoCronState.RepoOwner.ValueString()
	repoName := repoCronState.RepoName.ValueString()

	cron := prepareRepositoryCronPatch(repoCronPlan)

	cron, err := r.client.CronUpdate(repoOwner, repoName, cron)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository cron", err.Error())
		return
	}

	WoodpeckerToRepositoryCron(*cron, &repoCronState)

	diags = resp.State.Set(ctx, &repoCronState)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryCron) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var repoState RepositoryCron
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.RepoOwner.ValueString()
	repoName := repoState.RepoName.ValueString()
	repoId := repoState.ID.ValueInt64()

	err := r.client.CronDelete(repoOwner, repoName, repoId)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceRepositoryCron) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Need repo owner, repo name, and cron name?
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: repo_owner/repo_name/cron_name. Got: %s", req.ID),
		)
		return
	}

	// Search for cron to determine its ID

	resp.Diagnostics.AddError("todo: implement", "")
}

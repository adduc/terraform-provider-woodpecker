package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

func (r ResourceRepositoryCron) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Provides a repository resource. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/next/usage/cron)`,

		Attributes: map[string]schema.Attribute{
			// Required Attributes
			"repo_owner": schema.StringAttribute{
				Required:    true,
				Description: "User or organization responsible for repository",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repo_name": schema.StringAttribute{
				Required:    true,
				Description: "Repository name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Cron Name",
			},
			"schedule": schema.StringAttribute{
				Required:    true,
				Description: "Schedule (based on UTC)",
			},

			// Optional Attributes
			"repo_id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"creator_id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"next_exec": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"branch": schema.StringAttribute{
				Computed:    true,
				Description: "",
			},

			// Computed
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
			"created": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
		},
	}
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

	resourceData.RepoOwner = types.StringValue(repoOwner)
	resourceData.RepoName = types.StringValue(repoName)
	WoodpeckerToRepositoryCron(*cron, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryCron) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state RepositoryCron
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Created = state.Created
	plan.CreatorID = state.CreatorID
	plan.ID = state.ID

	if plan.RepoName.IsUnknown() {
		plan.RepoName = state.RepoName
	}

	if plan.RepoOwner.IsUnknown() {
		plan.RepoOwner = state.RepoOwner
	}

	if plan.Branch.IsUnknown() {
		plan.Branch = state.Branch
	}

	if plan.Name.IsUnknown() {
		plan.Name = state.Name
	}

	if plan.Schedule.IsUnknown() {
		plan.Schedule = state.Schedule
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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

	WoodpeckerToRepositoryCron(*cron, &repoCronPlan)

	diags = resp.State.Set(ctx, &repoCronPlan)
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
	idParts := strings.SplitN(req.ID, "/", 3)

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: repo_owner/repo_name/cron_name. Got: %s", req.ID),
		)
		return
	}

	repoOwner := idParts[0]
	repoName := idParts[1]
	cronName := idParts[2]

	crons, err := r.client.CronList(repoOwner, repoName)

	if err != nil {
		resp.Diagnostics.AddError("Could not fetch repository's cron list", err.Error())
		return
	}

	var cron RepositoryCron

	for _, wCron := range crons {
		if wCron.Name == cronName {
			WoodpeckerToRepositoryCron(*wCron, &cron)
			cron.RepoOwner = types.StringValue(repoOwner)
			cron.RepoName = types.StringValue(repoName)
			diags := resp.State.Set(ctx, &cron)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.AddError("Could not find cron with provided name", "")
}

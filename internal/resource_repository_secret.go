// A secret is composed of three identifiers:
// owner name, repository name, and secret name
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

func NewRepositorySecretResource() resource.Resource {
	return &ResourceRepositorySecret{}
}

type ResourceRepositorySecret struct {
	client woodpecker.Client
}

func (r ResourceRepositorySecret) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_secret"
}

func (r ResourceRepositorySecret) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Provides a repository secret. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/secrets)`,

		Attributes: map[string]tfsdk.Attribute{
			// Required Attributes
			"repo_owner": {
				Type:        types.StringType,
				Required:    true,
				Description: "User or organization responsible for repository",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"repo_name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Repository name",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "Secret Name",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"value": {
				Type:        types.StringType,
				Required:    true,
				Description: "Secret Value",
				Sensitive:   true,
			},

			// Optional Attributes
			"plugins_only": {
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				Description: "Whether secret is only available for plugins",
			},
			"images": {
				Type:        types.SetType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				Description: "List of images where this secret is available, leave empty to allow all images",
			},
			"events": {
				Type:        types.SetType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				Description: "One or more event types where secret is available (one of push, tag, pull_request, deployment, cron, manual)",
				Validators: []tfsdk.AttributeValidator{
					&ValidateSetInSlice{values: []string{"push", "tag", "pull_request", "deployment", "cron", "manual"}},
				},
			},

			// Computed
			"id": {
				Type:        types.Int64Type,
				Computed:    true,
				Description: "",
			},
		},
	}, nil
}

func (r *ResourceRepositorySecret) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceRepositorySecret) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositorySecret
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()

	secret, diags := prepareRepositorySecretPatch(ctx, resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cron, err := r.client.SecretCreate(repoOwner, repoName, secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not create repository cron", err.Error())
		return
	}

	resourceData.RepoOwner = types.String{Value: repoOwner}
	resourceData.RepoName = types.String{Value: repoName}
	diags = WoodpeckerToRepositorySecret(ctx, *cron, &resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositorySecret) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state RepositorySecret
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preknown attributes
	plan.ID = state.ID

	// Calculated / Configured

	if plan.RepoName.IsUnknown() {
		plan.RepoName = state.RepoName
	}

	if plan.RepoOwner.IsUnknown() {
		plan.RepoOwner = state.RepoOwner
	}

	if plan.Name.IsUnknown() {
		plan.Name = state.Name
	}

	if plan.Value.IsUnknown() {
		plan.Value = state.Value
	}

	if plan.PluginsOnly.IsUnknown() {
		plan.PluginsOnly = state.PluginsOnly
	}

	if plan.Images.IsUnknown() {
		plan.Images = state.Images
	}

	if plan.Events.IsUnknown() {
		plan.Events = state.Events
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositorySecret) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositorySecret
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()
	secretName := resourceData.Name.ValueString()

	secret, err := r.client.Secret(repoOwner, repoName, secretName)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToRepositorySecret(ctx, *secret, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositorySecret) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var repoSecretPlan RepositorySecret
	diags := req.Plan.Get(ctx, &repoSecretPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoSecretState RepositorySecret
	diags = req.State.Get(ctx, &repoSecretState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoSecretState.RepoOwner.ValueString()
	repoName := repoSecretState.RepoName.ValueString()

	secret, diags := prepareRepositorySecretPatch(ctx, repoSecretPlan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.SecretUpdate(repoOwner, repoName, secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository cron", err.Error())
		return
	}

	WoodpeckerToRepositorySecret(ctx, *secret, &repoSecretPlan)

	diags = resp.State.Set(ctx, &repoSecretPlan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositorySecret) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var repoState RepositorySecret
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.RepoOwner.ValueString()
	repoName := repoState.RepoName.ValueString()
	secretName := repoState.Name.ValueString()

	err := r.client.SecretDelete(repoOwner, repoName, secretName)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceRepositorySecret) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: repo_owner/repo_name/secret_name. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repo_owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repo_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
}

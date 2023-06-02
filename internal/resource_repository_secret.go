// A secret is composed of three identifiers:
// owner name, repository name, and secret name
package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

func (r ResourceRepositorySecret) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Provides a repository secret. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/secrets)`,

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
				Description: "Secret Name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "Secret Value",
				Sensitive:   true,
			},
			"events": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "One or more event types where secret is available (one of push, tag, pull_request, deployment, cron, manual)",

				Validators: []validator.Set{
					&ValidateSetInSlice{values: []string{"push", "tag", "pull_request", "deployment", "cron", "manual"}},
				},
			},

			// Optional Attributes
			"plugins_only": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether secret is only available for plugins",
			},
			"images": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of images where this secret is available, leave empty to allow all images",
			},

			// Computed
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
		},
	}
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
		resp.Diagnostics.AddError("Could not create repository secret", err.Error())
		return
	}

	resourceData.RepoOwner = types.StringValue(repoOwner)
	resourceData.RepoName = types.StringValue(repoName)
	diags = WoodpeckerToRepositorySecret(ctx, *cron, &resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositorySecret) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state RepositorySecret
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.Contains(plan.Name.ValueString(), "/") {
		resp.Diagnostics.AddError(
			"Unexpected character",
			"`/` is not supported in repository secret name",
		)
		return
	}

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

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
		resp.Diagnostics.AddError("Could not update repository secret", err.Error())
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
		resp.Diagnostics.AddError("Error deleting repository secret", err.Error())
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

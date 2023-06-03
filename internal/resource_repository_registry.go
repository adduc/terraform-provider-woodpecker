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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewRepositoryRegistryResource() resource.Resource {
	return &ResourceRepositoryRegistry{}
}

type ResourceRepositoryRegistry struct {
	client woodpecker.Client
}

func (r ResourceRepositoryRegistry) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_registry"
}

func (r ResourceRepositoryRegistry) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Provides a repository registry. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/registries)`,

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
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Registry Address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Registry Username",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "Registry Password",
				Sensitive:   true,
			},

			// Optional Attributes
			"token": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Registry Token",
				Sensitive:   true,
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Registry Email",
			},

			// Computed
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "",
			},
		},
	}
}

func (r *ResourceRepositoryRegistry) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceRepositoryRegistry) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryRegistry
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()

	registry, diags := prepareRepositoryRegistryPatch(ctx, resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registry, err := r.client.RegistryCreate(repoOwner, repoName, registry)

	if err != nil {
		resp.Diagnostics.AddError("Could not create repository registry", err.Error())
		return
	}

	resourceData.RepoOwner = types.StringValue(repoOwner)
	resourceData.RepoName = types.StringValue(repoName)
	diags = WoodpeckerToRepositoryRegistry(ctx, *registry, &resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryRegistry) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state RepositoryRegistry
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	if plan.Address.IsUnknown() {
		plan.Address = state.Address
	}

	if plan.Username.IsUnknown() {
		plan.Username = state.Username
	}

	if plan.Password.IsUnknown() {
		plan.Password = state.Password
	}

	if plan.Token.IsUnknown() {
		plan.Token = state.Token
	}

	if plan.Email.IsUnknown() {
		plan.Email = state.Email
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryRegistry) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData RepositoryRegistry
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	repoOwner := resourceData.RepoOwner.ValueString()
	repoName := resourceData.RepoName.ValueString()
	address := resourceData.Address.ValueString()

	registry, err := r.client.Registry(repoOwner, repoName, address)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToRepositoryRegistry(ctx, *registry, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryRegistry) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RepositoryRegistry
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state RepositoryRegistry
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := state.RepoOwner.ValueString()
	repoName := state.RepoName.ValueString()

	registry, diags := prepareRepositoryRegistryPatch(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registry, err := r.client.RegistryUpdate(repoOwner, repoName, registry)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository registry", err.Error())
		return
	}

	WoodpeckerToRepositoryRegistry(ctx, *registry, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceRepositoryRegistry) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var repoState RepositoryRegistry
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repoOwner := repoState.RepoOwner.ValueString()
	repoName := repoState.RepoName.ValueString()
	address := repoState.Address.ValueString()

	err := r.client.RegistryDelete(repoOwner, repoName, address)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository registry", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceRepositoryRegistry) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: repo_owner/repo_name/registry_address. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repo_owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repo_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("address"), idParts[2])...)
}

// A secret is composed of three identifiers:
// owner name, repository name, and secret name
package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewSecretResource() resource.Resource {
	return &ResourceSecret{}
}

type ResourceSecret struct {
	client woodpecker.Client
}

func (r ResourceSecret) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r ResourceSecret) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Provides a global secret. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/secrets).`,

		Attributes: map[string]tfsdk.Attribute{
			// Required Attributes
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

func (r *ResourceSecret) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceSecret) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData Secret
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, diags := prepareSecretPatch(ctx, resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cron, err := r.client.GlobalSecretCreate(secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not create repository secret", err.Error())
		return
	}

	diags = WoodpeckerToSecret(ctx, *cron, &resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceSecret) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state Secret
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

func (r ResourceSecret) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData Secret
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch secret
	secretName := resourceData.Name.ValueString()

	secret, err := r.client.GlobalSecret(secretName)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToSecret(ctx, *secret, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceSecret) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var repoSecretPlan Secret
	diags := req.Plan.Get(ctx, &repoSecretPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoSecretState Secret
	diags = req.State.Get(ctx, &repoSecretState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, diags := prepareSecretPatch(ctx, repoSecretPlan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.GlobalSecretUpdate(secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not update repository cron", err.Error())
		return
	}

	WoodpeckerToSecret(ctx, *secret, &repoSecretPlan)

	diags = resp.State.Set(ctx, &repoSecretPlan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceSecret) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var repoState Secret
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretName := repoState.Name.ValueString()

	err := r.client.GlobalSecretDelete(secretName)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting repository", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceSecret) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

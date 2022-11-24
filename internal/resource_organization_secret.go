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

func NewOrganizationSecretResource() resource.Resource {
	return &ResourceOrganizationSecret{}
}

type ResourceOrganizationSecret struct {
	client woodpecker.Client
}

func (r ResourceOrganizationSecret) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_secret"
}

func (r ResourceOrganizationSecret) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: `Provides a organization secret. For more 
		information see [Woodpecker CI's documentation](https://woodpecker-ci.org/docs/usage/secrets)`,

		Attributes: map[string]tfsdk.Attribute{
			// Required Attributes
			"owner": {
				Type:        types.StringType,
				Required:    true,
				Description: "Organization name",
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

func (r *ResourceOrganizationSecret) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceOrganizationSecret) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData OrganizationSecret
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	owner := resourceData.Owner.ValueString()

	secret, diags := prepareOrganizationSecretPatch(ctx, resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.OrgSecretCreate(owner, secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not create organization secret", err.Error())
		return
	}

	resourceData.Owner = types.StringValue(owner)
	diags = WoodpeckerToOrganizationSecret(ctx, *secret, &resourceData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceOrganizationSecret) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state OrganizationSecret
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.Contains(plan.Name.ValueString(), "/") {
		resp.Diagnostics.AddError(
			"Unexpected character",
			"`/` is not supported in organization secret name",
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

	if plan.Owner.IsUnknown() {
		plan.Owner = state.Owner
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

func (r ResourceOrganizationSecret) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData OrganizationSecret
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch repo
	owner := resourceData.Owner.ValueString()
	secretName := resourceData.Name.ValueString()

	secret, err := r.client.OrgSecret(owner, secretName)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToOrganizationSecret(ctx, *secret, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceOrganizationSecret) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationSecret
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoSecretState OrganizationSecret
	diags = req.State.Get(ctx, &repoSecretState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	owner := repoSecretState.Owner.ValueString()

	secret, diags := prepareOrganizationSecretPatch(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.OrgSecretUpdate(owner, secret)

	if err != nil {
		resp.Diagnostics.AddError("Could not update organization secret", err.Error())
		return
	}

	WoodpeckerToOrganizationSecret(ctx, *secret, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceOrganizationSecret) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var repoState OrganizationSecret
	diags := req.State.Get(ctx, &repoState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	owner := repoState.Owner.ValueString()
	secretName := repoState.Name.ValueString()

	err := r.client.OrgSecretDelete(owner, secretName)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting organization secret", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceOrganizationSecret) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected format: owner/secret_name. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}

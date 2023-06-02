package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func NewUserResource() resource.Resource {
	return &ResourceUser{}
}

type ResourceUser struct {
	client woodpecker.Client
}

func (r ResourceUser) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r ResourceUser) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a user resource.",

		Attributes: map[string]schema.Attribute{
			// Required Attributes
			"login": schema.StringAttribute{
				Required:    true,
				Description: "Username for user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Optional Attributes
			"email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Email address for user",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether user is active in the system",
			},

			// Computed Attributes
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "User ID",
			},
			"admin": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether user is a Woodpecker admin",
			},
			"avatar": schema.StringAttribute{
				Computed:    true,
				Description: "Avatar URL for user",
			},
		},
	}
}

func (r *ResourceUser) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ResourceUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// unmarshall request config into resourceData
	var resourceData User
	diags := req.Config.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch, diags := prepareUserPatch(ctx, resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UserPost(patch)
	if err != nil {
		resp.Diagnostics.AddError("Could not create user", err.Error())
		return
	}

	user, err := r.client.UserPatch(patch)
	if err != nil {
		resp.Diagnostics.AddError("Could not create user", err.Error())
		return
	}

	WoodpeckerToUser(ctx, *user, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceUser) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and recreate it
		return
	}

	var plan, state User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unchangeable Attributes
	plan.Login = state.Login

	// Computed Attributes
	plan.ID = state.ID
	plan.Admin = state.Admin
	plan.Avatar = state.Avatar

	// Optional Attributes
	if plan.Email.IsUnknown() {
		plan.Email = state.Email
	}

	if plan.Active.IsUnknown() {
		plan.Active = state.Active
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// unmarshall request config into resourceData
	var resourceData User
	diags := req.State.Get(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	login := resourceData.Login.ValueString()
	user, err := r.client.User(login)

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	WoodpeckerToUser(ctx, *user, &resourceData)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state User
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	patch, diags := prepareUserPatch(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo, err := r.client.UserPatch(patch)

	if err != nil {
		resp.Diagnostics.AddError("Could not update user", err.Error())
		return
	}

	WoodpeckerToUser(ctx, *repo, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r ResourceUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	login := state.Login.ValueString()
	err := r.client.UserDel(login)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r ResourceUser) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("login"), req.ID)...)
}

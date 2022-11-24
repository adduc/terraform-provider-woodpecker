package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func WoodpeckerToRepository(wRepo woodpecker.Repo, repo *Repository) {
	repo.ID = types.Int64Value(wRepo.ID)
	repo.Owner = types.StringValue(wRepo.Owner)
	repo.Name = types.StringValue(wRepo.Name)
	repo.FullName = types.StringValue(wRepo.FullName)
	repo.Avatar = types.StringValue(wRepo.Avatar)
	repo.Link = types.StringValue(wRepo.Link)
	repo.Kind = types.StringValue(wRepo.Kind)
	repo.Clone = types.StringValue(wRepo.Clone)
	repo.Branch = types.StringValue(wRepo.Branch)
	repo.Timeout = types.Int64Value(wRepo.Timeout)
	repo.Visibility = types.StringValue(wRepo.Visibility)
	repo.IsTrusted = types.BoolValue(wRepo.IsTrusted)
	repo.IsGated = types.BoolValue(wRepo.IsGated)
	repo.AllowPull = types.BoolValue(wRepo.AllowPull)
	repo.Config = types.StringValue(wRepo.Config)
}

func prepareRepositoryPatch(resourceData Repository) *woodpecker.RepoPatch {
	patch := woodpecker.RepoPatch{}

	if !resourceData.Config.IsNull() && !resourceData.Config.IsUnknown() {
		value := resourceData.Config.ValueString()
		patch.Config = &value
	}

	if !resourceData.IsTrusted.IsNull() && !resourceData.IsTrusted.IsUnknown() {
		value := resourceData.IsTrusted.ValueBool()
		patch.IsTrusted = &value
	}

	if !resourceData.IsGated.IsNull() && !resourceData.IsGated.IsUnknown() {
		value := resourceData.IsGated.ValueBool()
		patch.IsGated = &value
	}

	if !resourceData.Timeout.IsNull() && !resourceData.Timeout.IsUnknown() {
		value := resourceData.Timeout.ValueInt64()
		patch.Timeout = &value
	}

	if !resourceData.Visibility.IsNull() && !resourceData.Visibility.IsUnknown() {
		value := resourceData.Visibility.ValueString()
		patch.Visibility = &value
	}

	if !resourceData.AllowPull.IsNull() && !resourceData.AllowPull.IsUnknown() {
		value := resourceData.AllowPull.ValueBool()
		patch.AllowPull = &value
	}

	return &patch
}

func WoodpeckerToRepositoryCron(wCron woodpecker.Cron, cron *RepositoryCron) {
	cron.ID = types.Int64Value(wCron.ID)
	cron.Name = types.StringValue(wCron.Name)
	cron.RepoID = types.Int64Value(wCron.RepoID)
	cron.CreatorID = types.Int64Value(wCron.CreatorID)
	cron.NextExec = types.Int64Value(wCron.NextExec)
	cron.Schedule = types.StringValue(wCron.Schedule)
	cron.Created = types.Int64Value(wCron.Created)
	cron.Branch = types.StringValue(wCron.Branch)
}

func prepareRepositoryCronPatch(resourceData RepositoryCron) *woodpecker.Cron {
	patch := woodpecker.Cron{}

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Name.IsNull() && !resourceData.Name.IsUnknown() {
		patch.Name = resourceData.Name.ValueString()
	}

	if !resourceData.RepoID.IsNull() && !resourceData.RepoID.IsUnknown() {
		patch.RepoID = resourceData.RepoID.ValueInt64()
	}

	if !resourceData.CreatorID.IsNull() && !resourceData.CreatorID.IsUnknown() {
		patch.CreatorID = resourceData.CreatorID.ValueInt64()
	}

	if !resourceData.NextExec.IsNull() && !resourceData.NextExec.IsUnknown() {
		patch.NextExec = resourceData.NextExec.ValueInt64()
	}

	if !resourceData.Schedule.IsNull() && !resourceData.Schedule.IsUnknown() {
		patch.Schedule = resourceData.Schedule.ValueString()
	}

	if !resourceData.Created.IsNull() && !resourceData.Created.IsUnknown() {
		patch.Created = resourceData.Created.ValueInt64()
	}

	if !resourceData.Branch.IsNull() && !resourceData.Branch.IsUnknown() {
		patch.Branch = resourceData.Branch.ValueString()
	}

	return &patch
}

func WoodpeckerToRepositorySecret(ctx context.Context, wSecret woodpecker.Secret, secret *RepositorySecret) diag.Diagnostics {

	var diags, err diag.Diagnostics

	secret.ID = types.Int64Value(wSecret.ID)
	secret.Name = types.StringValue(wSecret.Name)

	if secret.Value.IsNull() || secret.Value.IsUnknown() {
		secret.Value = types.StringValue(wSecret.Value)
	}

	secret.PluginsOnly = types.BoolValue(wSecret.PluginsOnly)
	secret.Images, diags = types.SetValueFrom(ctx, types.StringType, wSecret.Images)
	secret.Events, err = types.SetValueFrom(ctx, types.StringType, wSecret.Events)

	diags.Append(err...)

	return diags
}

func prepareRepositorySecretPatch(ctx context.Context, resourceData RepositorySecret) (*woodpecker.Secret, diag.Diagnostics) {
	patch := woodpecker.Secret{}

	var diags, err diag.Diagnostics

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Name.IsNull() && !resourceData.Name.IsUnknown() {
		patch.Name = resourceData.Name.ValueString()
	}

	if !resourceData.Value.IsNull() && !resourceData.Value.IsUnknown() {
		patch.Value = resourceData.Value.ValueString()
	}

	if !resourceData.PluginsOnly.IsNull() && !resourceData.PluginsOnly.IsUnknown() {
		patch.PluginsOnly = resourceData.PluginsOnly.ValueBool()
	}

	if !resourceData.Images.IsNull() && !resourceData.Images.IsUnknown() {
		err = resourceData.Images.ElementsAs(ctx, &patch.Images, false)
		diags.Append(err...)
	}

	if !resourceData.Events.IsNull() && !resourceData.Events.IsUnknown() {
		err = resourceData.Events.ElementsAs(ctx, &patch.Events, false)
		diags.Append(err...)
	}

	return &patch, diags
}

func WoodpeckerToSecret(ctx context.Context, wSecret woodpecker.Secret, secret *Secret) diag.Diagnostics {

	var diags, err diag.Diagnostics

	secret.ID = types.Int64Value(wSecret.ID)
	secret.Name = types.StringValue(wSecret.Name)

	if secret.Value.IsNull() || secret.Value.IsUnknown() {
		secret.Value = types.StringValue(wSecret.Value)
	}

	secret.PluginsOnly = types.BoolValue(wSecret.PluginsOnly)
	secret.Images, diags = types.SetValueFrom(ctx, types.StringType, wSecret.Images)
	secret.Events, err = types.SetValueFrom(ctx, types.StringType, wSecret.Events)

	diags.Append(err...)

	return diags
}

func prepareSecretPatch(ctx context.Context, resourceData Secret) (*woodpecker.Secret, diag.Diagnostics) {
	patch := woodpecker.Secret{}

	var diags, err diag.Diagnostics

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Name.IsNull() && !resourceData.Name.IsUnknown() {
		patch.Name = resourceData.Name.ValueString()
	}

	if !resourceData.Value.IsNull() && !resourceData.Value.IsUnknown() {
		patch.Value = resourceData.Value.ValueString()
	}

	if !resourceData.PluginsOnly.IsNull() && !resourceData.PluginsOnly.IsUnknown() {
		patch.PluginsOnly = resourceData.PluginsOnly.ValueBool()
	}

	if !resourceData.Images.IsNull() && !resourceData.Images.IsUnknown() {
		err = resourceData.Images.ElementsAs(ctx, &patch.Images, false)
		diags.Append(err...)
	}

	if !resourceData.Events.IsNull() && !resourceData.Events.IsUnknown() {
		err = resourceData.Events.ElementsAs(ctx, &patch.Events, false)
		diags.Append(err...)
	}

	return &patch, diags
}

func WoodpeckerToUser(ctx context.Context, wUser woodpecker.User, user *User) {
	user.ID = types.Int64Value(wUser.ID)
	user.Login = types.StringValue(wUser.Login)
	user.Email = types.StringValue(wUser.Email)
	user.Avatar = types.StringValue(wUser.Avatar)
	user.Active = types.BoolValue(wUser.Active)
	user.Admin = types.BoolValue(wUser.Admin)
}

func prepareUserPatch(ctx context.Context, resourceData User) (*woodpecker.User, diag.Diagnostics) {
	patch := woodpecker.User{}

	var diags diag.Diagnostics

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Login.IsNull() && !resourceData.Login.IsUnknown() {
		patch.Login = resourceData.Login.ValueString()
	}

	if !resourceData.Email.IsNull() && !resourceData.Email.IsUnknown() {
		patch.Email = resourceData.Email.ValueString()
	}

	if !resourceData.Avatar.IsNull() && !resourceData.Avatar.IsUnknown() {
		patch.Avatar = resourceData.Avatar.ValueString()
	}

	if !resourceData.Active.IsNull() && !resourceData.Active.IsUnknown() {
		patch.Active = resourceData.Active.ValueBool()
	}

	if !resourceData.Admin.IsNull() && !resourceData.Admin.IsUnknown() {
		patch.Admin = resourceData.Admin.ValueBool()
	}

	return &patch, diags
}

func WoodpeckerToOrganizationSecret(ctx context.Context, wSecret woodpecker.Secret, secret *OrganizationSecret) diag.Diagnostics {

	var diags, err diag.Diagnostics

	secret.ID = types.Int64Value(wSecret.ID)
	secret.Name = types.StringValue(wSecret.Name)

	if secret.Value.IsNull() || secret.Value.IsUnknown() {
		secret.Value = types.StringValue(wSecret.Value)
	}

	secret.PluginsOnly = types.BoolValue(wSecret.PluginsOnly)
	secret.Images, diags = types.SetValueFrom(ctx, types.StringType, wSecret.Images)
	secret.Events, err = types.SetValueFrom(ctx, types.StringType, wSecret.Events)

	diags.Append(err...)

	return diags
}

func prepareOrganizationSecretPatch(ctx context.Context, resourceData OrganizationSecret) (*woodpecker.Secret, diag.Diagnostics) {
	patch := woodpecker.Secret{}

	var diags, err diag.Diagnostics

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Name.IsNull() && !resourceData.Name.IsUnknown() {
		patch.Name = resourceData.Name.ValueString()
	}

	if !resourceData.Value.IsNull() && !resourceData.Value.IsUnknown() {
		patch.Value = resourceData.Value.ValueString()
	}

	if !resourceData.PluginsOnly.IsNull() && !resourceData.PluginsOnly.IsUnknown() {
		patch.PluginsOnly = resourceData.PluginsOnly.ValueBool()
	}

	if !resourceData.Images.IsNull() && !resourceData.Images.IsUnknown() {
		err = resourceData.Images.ElementsAs(ctx, &patch.Images, false)
		diags.Append(err...)
	}

	if !resourceData.Events.IsNull() && !resourceData.Events.IsUnknown() {
		err = resourceData.Events.ElementsAs(ctx, &patch.Events, false)
		diags.Append(err...)
	}

	return &patch, diags
}

func WoodpeckerToRepositoryRegistry(ctx context.Context, wRegistry woodpecker.Registry, registry *RepositoryRegistry) diag.Diagnostics {

	var diags diag.Diagnostics

	registry.ID = types.Int64Value(wRegistry.ID)
	registry.Address = types.StringValue(wRegistry.Address)
	registry.Username = types.StringValue(wRegistry.Username)
	registry.Token = types.StringValue(wRegistry.Token)
	registry.Email = types.StringValue(wRegistry.Email)

	if registry.Password.IsNull() || registry.Password.IsUnknown() {
		registry.Password = types.StringValue(wRegistry.Password)
	}

	return diags
}

func prepareRepositoryRegistryPatch(ctx context.Context, resourceData RepositoryRegistry) (*woodpecker.Registry, diag.Diagnostics) {
	patch := woodpecker.Registry{}

	var diags diag.Diagnostics

	if !resourceData.ID.IsNull() && !resourceData.ID.IsUnknown() {
		patch.ID = resourceData.ID.ValueInt64()
	}

	if !resourceData.Address.IsNull() && !resourceData.Address.IsUnknown() {
		patch.Address = resourceData.Address.ValueString()
	}

	if !resourceData.Username.IsNull() && !resourceData.Username.IsUnknown() {
		patch.Username = resourceData.Username.ValueString()
	}

	if !resourceData.Password.IsNull() && !resourceData.Password.IsUnknown() {
		patch.Password = resourceData.Password.ValueString()
	}

	if !resourceData.Token.IsNull() && !resourceData.Token.IsUnknown() {
		patch.Token = resourceData.Token.ValueString()
	}

	if !resourceData.Email.IsNull() && !resourceData.Email.IsUnknown() {
		patch.Email = resourceData.Email.ValueString()
	}

	return &patch, diags
}

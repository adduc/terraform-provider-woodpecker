package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func WoodpeckerToRepository(wRepo woodpecker.Repo, repo *Repository) {
	repo.ID = types.Int64{Value: wRepo.ID}
	repo.Owner = types.String{Value: wRepo.Owner}
	repo.Name = types.String{Value: wRepo.Name}
	repo.FullName = types.String{Value: wRepo.FullName}
	repo.Avatar = types.String{Value: wRepo.Avatar}
	repo.Link = types.String{Value: wRepo.Link}
	repo.Kind = types.String{Value: wRepo.Kind}
	repo.Clone = types.String{Value: wRepo.Clone}
	repo.Branch = types.String{Value: wRepo.Branch}
	repo.Timeout = types.Int64{Value: wRepo.Timeout}
	repo.Visibility = types.String{Value: wRepo.Visibility}
	repo.IsTrusted = types.Bool{Value: wRepo.IsTrusted}
	repo.IsGated = types.Bool{Value: wRepo.IsGated}
	repo.AllowPull = types.Bool{Value: wRepo.AllowPull}
	repo.Config = types.String{Value: wRepo.Config}
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
	cron.ID = types.Int64{Value: wCron.ID}
	cron.Name = types.String{Value: wCron.Name}
	cron.RepoID = types.Int64{Value: wCron.RepoID}
	cron.CreatorID = types.Int64{Value: wCron.CreatorID}
	cron.NextExec = types.Int64{Value: wCron.NextExec}
	cron.Schedule = types.String{Value: wCron.Schedule}
	cron.Created = types.Int64{Value: wCron.Created}
	cron.Branch = types.String{Value: wCron.Branch}
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

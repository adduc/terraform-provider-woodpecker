package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
)

func WoodpeckerToRepository(wRepo woodpecker.Repo, repo *Repository) {
	repo.ID = types.Int64{Value: wRepo.ID}
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

	if !resourceData.Config.Null && !resourceData.Config.Unknown {
		patch.Config = &resourceData.Config.Value
	}

	if !resourceData.IsTrusted.Null && !resourceData.IsTrusted.Unknown {
		patch.IsTrusted = &resourceData.IsTrusted.Value
	}

	if !resourceData.IsGated.Null && !resourceData.IsGated.Unknown {
		patch.IsGated = &resourceData.IsGated.Value
	}

	if !resourceData.Timeout.Null && !resourceData.Timeout.Unknown {
		patch.Timeout = &resourceData.Timeout.Value
	}

	if !resourceData.Visibility.Null && !resourceData.Visibility.Unknown {
		patch.Visibility = &resourceData.Visibility.Value
	}

	if !resourceData.AllowPull.Null && !resourceData.AllowPull.Unknown {
		patch.AllowPull = &resourceData.AllowPull.Value
	}

	return &patch
}

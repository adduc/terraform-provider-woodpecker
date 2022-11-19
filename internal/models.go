package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Repository struct {
	ID         types.Int64  `tfsdk:"id"`
	Owner      types.String `tfsdk:"owner"`
	Name       types.String `tfsdk:"name"`
	FullName   types.String `tfsdk:"full_name"`
	Avatar     types.String `tfsdk:"avatar"`
	Link       types.String `tfsdk:"link"`
	Kind       types.String `tfsdk:"kind"`
	Clone      types.String `tfsdk:"clone"`
	Branch     types.String `tfsdk:"branch"`
	Timeout    types.Int64  `tfsdk:"timeout"`
	Visibility types.String `tfsdk:"visibility"`
	IsTrusted  types.Bool   `tfsdk:"is_trusted"`
	IsGated    types.Bool   `tfsdk:"is_gated"`
	AllowPull  types.Bool   `tfsdk:"allow_pull"`
	Config     types.String `tfsdk:"config"`
}

type Secret struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	PluginsOnly types.Bool   `tfsdk:"plugins_only"`
	Images      types.Set    `tfsdk:"images"`
	Events      types.Set    `tfsdk:"events"`
}

type User struct {
	ID     types.Int64  `tfsdk:"id"`
	Login  types.String `tfsdk:"login"`
	Email  types.String `tfsdk:"email"`
	Avatar types.String `tfsdk:"avatar"`
	Active types.Bool   `tfsdk:"active"`
	Admin  types.Bool   `tfsdk:"admin"`
}

type RepositoryCron struct {
	RepoOwner types.String `tfsdk:"repo_owner"`
	RepoName  types.String `tfsdk:"repo_name"`
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	RepoID    types.Int64  `tfsdk:"repo_id"`
	CreatorID types.Int64  `tfsdk:"creator_id"`
	NextExec  types.Int64  `tfsdk:"next_exec"`
	Schedule  types.String `tfsdk:"schedule"`
	Created   types.Int64  `tfsdk:"created"`
	Branch    types.String `tfsdk:"branch"`
}

type RepositorySecret struct {
	RepoOwner   types.String `tfsdk:"repo_owner"`
	RepoName    types.String `tfsdk:"repo_name"`
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	PluginsOnly types.Bool   `tfsdk:"plugins_only"`
	Images      types.Set    `tfsdk:"images"`
	Events      types.Set    `tfsdk:"events"`
}

type RepositorySecretData struct {
	RepoOwner   types.String `tfsdk:"repo_owner"`
	RepoName    types.String `tfsdk:"repo_name"`
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PluginsOnly types.Bool   `tfsdk:"plugins_only"`
	Images      types.Set    `tfsdk:"images"`
	Events      types.Set    `tfsdk:"events"`
}

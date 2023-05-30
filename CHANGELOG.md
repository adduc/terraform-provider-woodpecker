# Changelog

## v0.4.0 - 2023-XX-XX

### Added

- Add tests for common use cases

### Changed

- Upgrade to Terraform plugin framework 0.17.0

## v0.3.0 - 2022-11-24

### Added

- Add resource `woodpecker_repository_registry`
- Add data-source `woodpecker_repository_registry`
- Add resource `woodpecker_organization_secret`
- Add data-source `woodpecker_organization_secret`
- Add resource `woodpecker_user`
- Add data-source `woodpecker_user`

### Changed

- secret: an error is now triggered when attempting to create a secret
  or repository secret with a slash (`/`) in its name. While Woodpecker
  creates secrets successfully, any attempt to update or delete a secret
  that has a slash in its name silently fails.

### Fixed

- Fix inability to import crons with slashes (`/`) in their name.

## v0.2.0 - 2022-11-20

### Added

- Add resource `woodpecker_secret`
- Add data-source `woodpecker_secret`

### Changed

- Update documentation (added examples, provider schema documentation)
- Upgrade to Terraform plugin framework 0.16.0

### Fixed

- Fix error messages using wrong resource name

## v0.1.0 - 2022-11-13

### Added

- Add `repository_secret` resource and data sources. (#8, #9)

### Fixed

- Repository: the list of repositories is now refreshed when adding a new repository to fix an issue with newly created repositories not being found by the provider.
- Fixed an issue where changes to repository owner or repository name were not included in the planned list of changes.
- Resources are now properly marked for recreation when dependent fields are changed (repository owner, repository name, secret name, etc.).
- Update documentation

## v0.0.1 - 2022-11-09

Initial release, with support for administering repository and
repository cron data sources and resources.

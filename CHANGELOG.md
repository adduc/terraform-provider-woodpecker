# Changelog

## [Unreleased]

### Changed

- secret: an error is now triggered when attempting to create a secret
  or repository secret with a slash (`/`) in its name. While Woodpecker
  creates secrets successfully, any attempt to update or delete a secret
  that has a slash in its name silently fails.

### Fixed

- Fix inability to import crons with slashes (`/`) in their name.


## v0.2.0 - 2022-11-20

### Added

- Added resource `woodpecker_secret`
- Added data-source `woodpecker_secret`

### Changed

- Updated documentation (added examples, provider schema documentation)
- Upgraded to Terraform plugin framework 0.16.0

### Fixed

- Fixed error messages using wrong resource name


## v0.1.0 - 2022-11-13

### Added

* Added `repository_secret` resource and data sources. (#8, #9)

### Fixed

* Repository: the list of repositories is now refreshed when adding a new repository to fix an issue with newly created repositories not being found by the provider.
* Fixed an issue where changes to repository owner or repository name were not included in the planned list of changes.
* Resources are now properly marked for recreation when dependent fields are changed (repository owner, repository name, secret name, etc.).
* Updated documentation


## v0.0.1 - 2022-11-09

Initial release, with support for administering repository and
repository cron data sources and resources.

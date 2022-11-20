# Changelog

## [Unreleased]

### Changed

- secret: an error is now triggered when attempting to create a secret
  or repository secret with a slash (`/`) in its name. While Woodpecker
  creates secrets successfully, any attempt to update or delete a secret
  that has a slash in its name silently fails.

### Fixed

- Fix inability to import crons with slashes (`/`) in their name.
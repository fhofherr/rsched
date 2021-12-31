# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

* Binaries are built with linker flag `-s`. This creates a smaller
  binary.

## [0.1.0-alpha.3] - 2021-12-28

### Changed

* Environment variables used to configure `rsched` now need to be
  prefixed with `RSCHED_` instead of `RSCHED__`.
* Docker image is now based on Alpine Linux.
* Fix injecting version information into binaries.
* Copy `restic` and `rsched` binaries to `/usr/local/bin` when creating
  image.

## [0.1.0-alpha.2] - 2021-12-26

* Add additional log statements to make debugging easier.

## [0.1.0-alpha.1] - 2021-12-26

### Added

* Ability to create backups usign [`restic`](https://restic.net/). In
  order to create a one of backup the special value `once` may be passed
  to the `-restic-schedule` flag.
* `Makefile` that helps with building and allows to create and publish
  docker images. Currently only `linux/amd64` is supported.
  `linux/arm64` will be supported once I have the ability to test it.

[Unreleased]: https://github.com/fhofherr/rsched/compare/v0.1.0-alpha.3...HEAD
[0.1.0-alpha.3]: https://github.com/fhofherr/rsched/compare/v0.1.0-alpha.2...v0.1.0-alpha.3
[0.1.0-alpha.2]: https://github.com/fhofherr/rsched/compare/v0.1.0-alpha.1...v0.1.0-alpha.2
[0.1.0-alpha.1]: https://github.com/fhofherr/rsched/compare/v0.0.0...v0.1.0-alpha.1
[0.0.0]: https://github.com/fhofherr/rsched/releases/tag/v0.0.0

# Change Log

**ATTN**: This project uses [semantic versioning](http://semver.org/).

## [Unreleased]

### Added

### Changed

## [0.2.0] - unreleased
### Added
- Support for variadic handlers in `New()`
- Added `Negroni.Handlers()` to fetch all of the handlers for a given chain
- Allowed size in `Recovery` handler was bumped to 8k
- `Negroni.UseFunc` to push another handler onto the chain

### Changed
- Set the status before calling `beforeFuncs` so the information is available to them
- Set default status to `200` in the case that no handler writes status -- was previously `0`

## 0.1.0 - 2013-07-22
### Added
- Initial implementation.

[Unreleased]: https://github.com/codegangsta/negroni/compare/v0.1.0...HEAD

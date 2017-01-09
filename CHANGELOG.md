# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) 
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.1.0] - 2017-01-09

### Changed

- Use [`qingstor-sdk`](https://github.com/yunify/qingstor-sdk-python) instead
- Default config path changed to `~/.qingstor`

### Fixed

- Fix bug when listing buckets under python 3
- Return -1 while download failed

### BREAKING CHANGE

- Config should be updated to new version, older version will be no more supported.

## [1.0.5] - 2016-11-29

### Added

- Catch the exception for file not found

## [1.0.4] - 2016-09-11

### Added

- Add `.gitignore`
- Add MIME Type Detection on File Upload

### Fixed

- Fix bug when printing help pages
- Fix bug when removing empty local directories

## [1.0.3] - 2016-07-26

### Changed

- `PART_SIZE` changed to 32MB
- Validate bucket name before command performs

## 1.0.0 - 2016-07-05

### Added

- Hello, qsctl.

[1.1.0]: https://github.com/yunify/qsctl/compare/2cc5fe3c912dc37356d332b103c0f132e1058c63...1.1.0
[1.0.5]: https://github.com/yunify/qsctl/compare/82d42dcaaec9d58c8fdd6cad82bac56092416ff6...2cc5fe3c912dc37356d332b103c0f132e1058c63
[1.0.4]: https://github.com/yunify/qsctl/compare/3073a03e7d2d801226c525e574f9bba295e12ddd...82d42dcaaec9d58c8fdd6cad82bac56092416ff6
[1.0.3]: https://github.com/yunify/qsctl/compare/69a52585edb6b14247e8954722d7b6e680769612...3073a03e7d2d801226c525e574f9bba295e12ddd

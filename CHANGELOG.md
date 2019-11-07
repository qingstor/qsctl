# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.0.0-beta.1] - 2019-11-07

### Added

- Rm -r based on async list (#132)
- Implement feature rb -f (#134)
- makefile: Make go build reproducible (#148)
- cmd: Support set config interactively (#152)
- task: Implement workload based scheduler (#159)
- cmd/cp: Implement cp -r (#155)
- cmd/sync: Implement sync (#162)
- cmd/cp: Add support for cp -r from root path. (#164)

### Changed

- Replace local file system operations to posixfs storager (#130)
- *: Refactor task scheduler system (#146)
- pkg/schedule: Submit task directly (#147)
- *: Refactor task type system (#151)

### Fixed

- cmd/ls: fix ls output, exit until output finish (#128)
- task: Fix segment upload not correctly (#154)
- cmd/utils/setup: Fix bug while creating config file but directory not exists (#165)

## [2.0.0-alpha.8] - 2019-10-17

### Added

- cmd/stat: Add updatedAt output support (#123)
- cmd/ls: Add ls command, use stream stdout while list objects (#122, #121, #120)
- Add fault support (#113)
- utils: Add bucket name validator (#109)
- misc: Add support for codecov (#108)

### Changed

- Replace storage logic with new storage layer (#118)
- cmd: Unify error handle by using TriggerFault (#116)
- pkg: Move types to pkg and merge utils for better share (#111)
- Use new task execution framework (#107)
- Modify cmd by using new task framework: cp(#107), mb(#110), presign(#114), rb(#115), rm(#119), stat(#112)

## [2.0.0-alpha.7] - 2019-07-15

### Added

- action/stat: Add format support for stat (#95)

### Fixed

- cmd/stat: Use fmt.Println instead to print to stdout (#94)

## [2.0.0-alpha.6] - 2019-07-14

### Added

- cmd: Add rb command to remove an empty bucket (#91)

### Changed

- main: Refactor config file load logic (#92)

## [2.0.0-alpha.5] - 2019-07-11

### Added

- action: Add unit tests for copy (#88)

### Fixed

- action/copy: Fix wrong path parsed in copying from remote (#87)

## [2.0.0-alpha.4] - 2019-07-11

### Changed

-  misc: Upgrade go-sdk to 3.0.1

### Fixed

- action/copy: Fix md5 sum incorrectly (#85)

## [2.0.0-alpha.3] - 2019-07-10

### Added

- action: Add unit tests for utils
- Add mb make bucket command
- storage/mock: Implement MockObjectStorage (#82)

### Changed

- action/utils: Refactor ParseQsPath
- Refactor the way to init commands' flags and global service
- action/bucket: Unify the format of logs, and modify bucket created callback into log
- action/bucket: Move put bucket logic by qingstor into helpers
- cmd/mb: modify the way to check required flags (#78)
- storage: Convert helper into storage interface for better testing (#79)


### Fixed

- cmd: Fix subcommands' flags been polluted wrongly
- action/utils: Fix object key not parsed correctly
- action/utils: Fix bucket name can't starts or ends with "-"
- action/copy: Fix wg.Add(1) not executed after task submit
- cmd/rm: Implement rm command for single object (#81)

## [2.0.0-alpha.2] - 2019-07-08

### Fixed

- action/copy: Fix buffer writer not flushed
- action/utils: "-" in bucket name should be allowed

## [2.0.0-alpha.1] - 2019-07-05

### Added

- Hello, qsctl, again.

## [1.7.7] - 2018-12-01

### Changed

- Improve the logic for list-objects paging handling
- Add support for uploading empty file

## [1.7.6] - 2018-08-21

### Fixed

- Fix bug which leads to wrong key deleted

## [1.7.5] - 2017-12-15

### Fixed

- Fix memory leak in task queue
- Fix python2 on x86 can't convert int to long
- Fix si_unit is not callable
- Fix encoding error with python 2 on windows platform

## [1.7.4] - 2017-10-03

### Changed

- Use larger BUFFER_SIZE for downloads

### Fixed

- Fix bug lead to mkdir conflict in multi threads

## [1.7.3] - 2017-09-28

### Added

- Add SectionReader to replace StringIO

### Changed

- Use 64M for part size
- Fixed read blocksize to 4M

### Fixed

- Improve upload speed on windows platform

## [1.7.2] - 2017-09-20

### Fixed

- Broken dependence on python 2

## [1.7.1] - 2017-09-14

### Added

- Add zone options support
- Add tab completion for linux

### Changed

- Handle illegal characters in a better way
- Modify current_bucket to variable of class

### Fixed

- Fix arg conflict bug
- Fix bug which leads to recursion depth exceeded

## [1.7.0] - 2017-08-11

### Added

- Add support for public bucket
- Add multi thread support for upload and download

### Changed

- Refactor presign command
- Use independent threads to process output
- Use a separate thread to process the progress bar
- Refactor convert_to_bytes function
- Refactor UploadIDRecoder class

## [1.6.2] - 2017-06-20

### Fixed

- Zone should be required while create bucket

## [1.6.1] - 2017-05-23

### Changed

- Use class variable to pass options

### Fixed

- Fix SSL Warnings with old python versions
- Fix part_numbers is not defined error
- Fix bug that cause presign can't work

## [1.6.0] - 2017-05-13

### Added

- Add rate limit feature

### Fixed

- Fix crash when executing the help sub-command

## [1.5.0] - 2017-04-11

### Added

- Add resumimg downloading at break-point
- Add resumimg uploading at break-point
- Handle interrupt signal silently

### Fixed

- Fix duplicate output while detecting config file
- Fix pkg_resource error in python 2.6
- Fix cross-platform coding problems

## [1.4.1] - 2017-03-25

### Fixed

- Fix tqdm encoding error on python 2.6

## [1.4.0] - 2017-03-01

### Added

- Add presign command

### Changed

- Allow user operating buckets granted by policy

### Fixed

- Fix progressbar not close correctly

## [1.3.1] - 2017-02-28

### Fixed

- Fix bug in sync command

## [1.3.0] - 2017-02-27

### Added

- Add progress bar while downloading and uploading

### Changed

- Use DeleteMultipleObjects API instead

### Fixed

- Fix bug while deleting not exist object
- Fix force argument's wrong behavior on multipart
- Fix confirm statement encoding error in python2

## [1.2.3] - 2017-02-08

### Fixed

- Fix bug in using old version config

## [1.2.2] - 2017-01-20

### Changed

- Refactor config file load function, support load config from `~/.qingcloud`
- Be compatible with `qy_access_key_id` and `qy_secret_access_key`

### Fixed

- Fix bug while params is int instead of str

## [1.2.1] - 2017-01-12

### Fixed

- Fix import error in python3

## [1.2.0] - 2017-01-10

### Added

- Support to upload file from stdin

### Fixed

- Fix bug that list_keys do not respect prefix
- Fix empty output while access_key_id is invalid

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

[2.0.0-beta.1]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.8...v2.0.0-beta.1
[2.0.0-alpha.8]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.7...v2.0.0-alpha.8
[2.0.0-alpha.7]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.6...v2.0.0-alpha.7
[2.0.0-alpha.6]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.5...v2.0.0-alpha.6
[2.0.0-alpha.5]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.4...v2.0.0-alpha.5
[2.0.0-alpha.4]: https://github.com/yunify/qsctl/compare/v2.0.0-alpha.3...v2.0.0-alpha.4
[2.0.0-alpha.3]: https://github.com/yunify/qsctl/compare/2.0.0-alpha.2...2.0.0-alpha.3
[2.0.0-alpha.2]: https://github.com/yunify/qsctl/compare/2.0.0-alpha.1...2.0.0-alpha.2
[2.0.0-alpha.1]: https://github.com/yunify/qsctl/compare/1.7.7...2.0.0-alpha.1
[1.7.7]: https://github.com/yunify/qsctl/compare/1.7.6...1.7.7
[1.7.6]: https://github.com/yunify/qsctl/compare/1.7.5...1.7.6
[1.7.5]: https://github.com/yunify/qsctl/compare/1.7.4...1.7.5
[1.7.4]: https://github.com/yunify/qsctl/compare/1.7.3...1.7.4
[1.7.3]: https://github.com/yunify/qsctl/compare/1.7.2...1.7.3
[1.7.2]: https://github.com/yunify/qsctl/compare/1.7.1...1.7.2
[1.7.1]: https://github.com/yunify/qsctl/compare/1.7.0...1.7.1
[1.7.0]: https://github.com/yunify/qsctl/compare/1.6.2...1.7.0
[1.6.2]: https://github.com/yunify/qsctl/compare/1.6.1...1.6.2
[1.6.1]: https://github.com/yunify/qsctl/compare/1.6.0...1.6.1
[1.6.0]: https://github.com/yunify/qsctl/compare/1.5.0...1.6.0
[1.5.0]: https://github.com/yunify/qsctl/compare/1.4.1...1.5.0
[1.4.1]: https://github.com/yunify/qsctl/compare/1.4.0...1.4.1
[1.4.0]: https://github.com/yunify/qsctl/compare/1.3.1...1.4.0
[1.3.1]: https://github.com/yunify/qsctl/compare/1.3.0...1.3.1
[1.3.0]: https://github.com/yunify/qsctl/compare/1.2.3...1.3.0
[1.2.3]: https://github.com/yunify/qsctl/compare/1.2.2...1.2.3
[1.2.2]: https://github.com/yunify/qsctl/compare/1.2.1...1.2.2
[1.2.1]: https://github.com/yunify/qsctl/compare/1.2.0...1.2.1
[1.2.0]: https://github.com/yunify/qsctl/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/yunify/qsctl/compare/2cc5fe3c912dc37356d332b103c0f132e1058c63...1.1.0
[1.0.5]: https://github.com/yunify/qsctl/compare/82d42dcaaec9d58c8fdd6cad82bac56092416ff6...2cc5fe3c912dc37356d332b103c0f132e1058c63
[1.0.4]: https://github.com/yunify/qsctl/compare/3073a03e7d2d801226c525e574f9bba295e12ddd...82d42dcaaec9d58c8fdd6cad82bac56092416ff6
[1.0.3]: https://github.com/yunify/qsctl/compare/69a52585edb6b14247e8954722d7b6e680769612...3073a03e7d2d801226c525e574f9bba295e12ddd

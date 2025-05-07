# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.4.9] - 2025-05-07

### Added

- goreleaser: Add goreleaser to build and publish binary (#398)

## [2.4.8] - 2025-04-27

### Added

- feat: feat: Allow object keys include multiple slashes if configured (#396)

## [2.4.7] - 2023-08-21

### Fixes

- fix(deps): Update qingstor lib to fix listDir behavior (#392)

## [2.4.6] - 2023-07-18

### Added

- feat: Added support for virtual host style (#390)


## [2.4.5] - 2022-04-19

### Changed

- ci: Update go to recent two versions (#379)

## [2.4.4] - 2022-01-06

### Fixed

- cmd/tee: Fixed the part exceeds maximum number problem (#376)

## [2.4.3] - 2021-05-12

### Changed

- build: Use static build via cgo disabled (#371)

## [2.4.2] - 2021-02-02

### Changed

- ci: Remove travis, add github action workflows (#353)

### Fixed

- cmd: Fix exit code not 1 when exit with error (#367)

## [2.4.1] - 2020-10-30

### Fixed

- cmd/mv: Fix mv command failed with valid check (#348)

## [2.4.0] - 2020-09-21

### Added

- cmd: Support include-regexp and exclude-regexp in sync command (#343)
- cmd: Support custom part size and threshold for multipart upload (#341)

### Fixed

- taskutils: Fix progress handler not wait as expected (#340)

## [2.3.0] - 2020-08-19

### Added

- cmd/shell: Add support for progress bar in transfer commands (#331)
- cmd: Support context in commands, support interrupt command in shell (#330)
- cmd/shell: Add exit command to exit shell

### Fixed

- cmd/shell: Modify initialing bucket list asynchronously to fix shell stuck when init (#322)
- cmd/stat: Fixed qsctl panic when execute stat command

## [2.2.0] - 2020-07-30

### Added

- cmd/shell: Add support for interactive shell (#310)

### Fixed

-  cmd/rb: Fix zone flag did not work (#314)

### Removed

- Remove interactive action in commands except shell

## [2.1.2] - 2020-06-30

### Changed

- mod: Bump storage version to adapt system specific separator to slash (#302)

### Fixed

- cmd/progress: Fix progress panic with negative WaitGroup counter (#296)

## [2.1.1] - 2020-06-03

### Changed

- i18n: Set language to en-US if detect failed instead of exit 1 (#291)

### Fixed

- cmd/utils: Fix qsctl hang on non-interactive environment (#292)
- Fix errors not return correctly in task exec

## [2.1.0] - 2020-05-28

### Added

- cmd/stat: Add support to stat a bucket (#279)
- cmd: Add global zone flag to assign zone manually (#282)
- cmd/ls: Add support to ls buckets by long format (#280)

### Fixed

- cmd/progress-bar: Fix data race in progress bar (#284)

## [2.0.1] - 2020-05-21

### Changed

- misc: Modify install path to /usr/bin in rpm/deb. (#275)

### Fixed

- mod: Upgrade noah to fix cp large file from qingstor. (#276)

## [2.0.0] - 2020-05-12

### Added

- cmd: Add support for no-progress flag (#262)
- *: Add support to make package for linux distribution (#263)

## [2.0.0-rc.2] - 2020-05-08

### Added

- cmd: Add double check when rb -f and confirm check when rm (#253)
- cmd: Improve the performance of progress bar (#256)

## [2.0.0-rc.1] - 2020-04-08

### Added

- cmd/sync: Add more flags support for sync. (#221)

### Changed

- *: Split task to qingstor/noah (#195)
- cmd: Silence usage when handled error returns. (#206)
- cmd/stat: Modify stat file name, use full name instead. (#209)
- cmd/ls: Align ls output while -h flag was set. (#248)

### Fixed

- utils: Fix bug that qs work dir not work in Windows. (#212)
- cmd/stat: Fix md5 not return in stat. (#216)
- cmd/ls: Fix ls behave differently with recursive flag (#219)

## [2.0.0-beta.2] - 2019-12-10

### Added

- cmd/mv: Implement feature mv (#171)
- cmd/sync: Implement ignore existing while sync file (#172)
- i18n: Add whole i18n support (#184)

### Changed

- *: Only print debug log while user required (#180)
- *: Move project to QingStor
- cmd: Improve qsctl interactive experience (#191)

### Fixed

- cmd: Fix object size not used correctly (#189)

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

[2.4.8]: https://github.com/qingstor/qsctl/compare/v2.4.8...v2.4.9
[2.4.8]: https://github.com/qingstor/qsctl/compare/v2.4.7...v2.4.8
[2.4.7]: https://github.com/qingstor/qsctl/compare/v2.4.6...v2.4.7
[2.4.6]: https://github.com/qingstor/qsctl/compare/v2.4.5...v2.4.6
[2.4.5]: https://github.com/qingstor/qsctl/compare/v2.4.4...v2.4.5
[2.4.4]: https://github.com/qingstor/qsctl/compare/v2.4.3...v2.4.4
[2.4.3]: https://github.com/qingstor/qsctl/compare/v2.4.2...v2.4.3
[2.4.2]: https://github.com/qingstor/qsctl/compare/v2.4.1...v2.4.2
[2.4.1]: https://github.com/qingstor/qsctl/compare/v2.4.0...v2.4.1
[2.4.0]: https://github.com/qingstor/qsctl/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/qingstor/qsctl/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/qingstor/qsctl/compare/v2.1.2...v2.2.0
[2.1.2]: https://github.com/qingstor/qsctl/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/qingstor/qsctl/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/qingstor/qsctl/compare/v2.0.1...v2.1.0
[2.0.1]: https://github.com/qingstor/qsctl/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/qingstor/qsctl/compare/v2.0.0-rc.2...v2.0.0
[2.0.0-rc.2]: https://github.com/qingstor/qsctl/compare/v2.0.0-rc.1...v2.0.0-rc.2
[2.0.0-rc.1]: https://github.com/qingstor/qsctl/compare/v2.0.0-beta.2...v2.0.0-rc.1
[2.0.0-beta.2]: https://github.com/qingstor/qsctl/compare/v2.0.0-beta.1...v2.0.0-beta.2
[2.0.0-beta.1]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.8...v2.0.0-beta.1
[2.0.0-alpha.8]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.7...v2.0.0-alpha.8
[2.0.0-alpha.7]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.6...v2.0.0-alpha.7
[2.0.0-alpha.6]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.5...v2.0.0-alpha.6
[2.0.0-alpha.5]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.4...v2.0.0-alpha.5
[2.0.0-alpha.4]: https://github.com/qingstor/qsctl/compare/v2.0.0-alpha.3...v2.0.0-alpha.4
[2.0.0-alpha.3]: https://github.com/qingstor/qsctl/compare/2.0.0-alpha.2...2.0.0-alpha.3
[2.0.0-alpha.2]: https://github.com/qingstor/qsctl/compare/2.0.0-alpha.1...2.0.0-alpha.2
[2.0.0-alpha.1]: https://github.com/qingstor/qsctl/compare/1.7.7...2.0.0-alpha.1
[1.7.7]: https://github.com/qingstor/qsctl/compare/1.7.6...1.7.7
[1.7.6]: https://github.com/qingstor/qsctl/compare/1.7.5...1.7.6
[1.7.5]: https://github.com/qingstor/qsctl/compare/1.7.4...1.7.5
[1.7.4]: https://github.com/qingstor/qsctl/compare/1.7.3...1.7.4
[1.7.3]: https://github.com/qingstor/qsctl/compare/1.7.2...1.7.3
[1.7.2]: https://github.com/qingstor/qsctl/compare/1.7.1...1.7.2
[1.7.1]: https://github.com/qingstor/qsctl/compare/1.7.0...1.7.1
[1.7.0]: https://github.com/qingstor/qsctl/compare/1.6.2...1.7.0
[1.6.2]: https://github.com/qingstor/qsctl/compare/1.6.1...1.6.2
[1.6.1]: https://github.com/qingstor/qsctl/compare/1.6.0...1.6.1
[1.6.0]: https://github.com/qingstor/qsctl/compare/1.5.0...1.6.0
[1.5.0]: https://github.com/qingstor/qsctl/compare/1.4.1...1.5.0
[1.4.1]: https://github.com/qingstor/qsctl/compare/1.4.0...1.4.1
[1.4.0]: https://github.com/qingstor/qsctl/compare/1.3.1...1.4.0
[1.3.1]: https://github.com/qingstor/qsctl/compare/1.3.0...1.3.1
[1.3.0]: https://github.com/qingstor/qsctl/compare/1.2.3...1.3.0
[1.2.3]: https://github.com/qingstor/qsctl/compare/1.2.2...1.2.3
[1.2.2]: https://github.com/qingstor/qsctl/compare/1.2.1...1.2.2
[1.2.1]: https://github.com/qingstor/qsctl/compare/1.2.0...1.2.1
[1.2.0]: https://github.com/qingstor/qsctl/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/qingstor/qsctl/compare/2cc5fe3c912dc37356d332b103c0f132e1058c63...1.1.0
[1.0.5]: https://github.com/qingstor/qsctl/compare/82d42dcaaec9d58c8fdd6cad82bac56092416ff6...2cc5fe3c912dc37356d332b103c0f132e1058c63
[1.0.4]: https://github.com/qingstor/qsctl/compare/3073a03e7d2d801226c525e574f9bba295e12ddd...82d42dcaaec9d58c8fdd6cad82bac56092416ff6
[1.0.3]: https://github.com/qingstor/qsctl/compare/69a52585edb6b14247e8954722d7b6e680769612...3073a03e7d2d801226c525e574f9bba295e12ddd

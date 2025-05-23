# qsctl

[![Build Status](https://github.com/qingstor/qsctl/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/qingstor/qsctl/actions?query=workflow%3A%22Unit+Test%22)
[![GoDoc](https://godoc.org/github.com/qingstor/qsctl?status.svg)](https://godoc.org/github.com/qingstor/qsctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/qingstor/qsctl)](https://goreportcard.com/report/github.com/qingstor/qsctl)
[![codecov](https://codecov.io/gh/qingstor/qsctl/branch/master/graph/badge.svg)](https://codecov.io/gh/qingstor/qsctl)
[![localized](https://badges.crowdin.net/qsctl/localized.svg)](https://crowdin.com/project/qsctl)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/qingstor/qsctl/blob/master/LICENSE)
[![Join the Chat](https://img.shields.io/badge/chat-online-blue?style=flat&logo=zulip)](https://qingstor.zulipchat.com/join/wiqqtnvhnux73i3q2qvep7ak/)

[中文文档](./docs/README-zh_CN.md)

qsctl is intended to be an advanced command line tool for QingStor, it provides
powerful unix-like commands to let you manage QingStor resources just like files
on local machine.

## Installation

### Binary

Visit <https://github.com/qingstor/qsctl/releases> to get latest releases.

## Getting Started

### Configure via initialization wizard

### Configure manually

To use qsctl, there must be a configuration file , for example

```yaml
access_key_id: 'ACCESS_KEY_ID_EXAMPLE'
secret_access_key: 'SECRET_ACCESS_KEY_EXAMPLE'
```

The configuration file is `~/.qingstor/config.yaml` by default, it also
can be specified by the option `-c /path/to/config`.

You can also config other option like `host` , `port` and so on, just
add lines below into configuration file, for example

```yaml
host: 'qingstor.com'
port: 443
protocol: 'https'
connection_retries: 3
zone: 'zone_name'

# By default, extra slashes in object key will be stripped from the path,
# like `/a//b` become `a/b`, this option disable that behavior.
# disable_uri_cleaning: true
```

~~You can also run qsctl command without config file, it will try to
find config file in specific directories, and if none of them contain
a config file, there will be an interactive setup to help you create
the config file, which will be created at `{$HOME}/.qingstor/config.yaml`.
(PS: The specific config file path depends on your os, usually
`~/.qingstor/config.yaml` in unix-like os, and
`C:\User\{username}\.qingstor\config.yaml` in Windows.)~~
```
interactive setup was removed from v2.2.0, and may be added in the future with an independent command.
```

Since v2.2.0, we moved all interactive operation into `qsctl shell`, in order that you can call commands in your
script without any interactive interruption.
We also removed progress bar rendering from all commands and added it into `qsctl shell`.

## Qsctl Shell

We introduced interactive shell from v2.2.0, which contains more instruction and tips. It is highly recommended
for those who just start to use qsctl. You can just execute `qsctl shell` to enter the shell interface and run commands
according to the tips. We also support commands history and show tips and auto completion for:
- available commands at the beginning of new line.
- buckets after inputting `qs://`.
- local files after inputting space.
- available flags for current command.

## Available Commands

Commands supported by qsctl are listed below:

- `cat`: Cat a remote object into stdout.
- `cp`: Copy local file(s) to QingStor or QingStor key(s) to local.
- `ls`: List buckets, or objects with given prefix.
- `mb`: Make a new bucket.
- `mv`: Move local file(s) to QingStor or QingStor key(s) to local
- `presign`: Get the pre-signed URL by given object key.
- `rb`: Delete a bucket.
- `rm`: Remove remote object(s).
- `shell`: start an interactive shell of qsctl.
- `stat`: Stat a remote object.
- `sync`: Sync between local directory and QS-Directory.
- `tee`: Tee from stdin to a remote object. (NOTICE: qsctl will not tee the content to stdout like linux tee command does.)

See the detailed usage and examples with `qsctl --help` or `qsctl <command> --help`.

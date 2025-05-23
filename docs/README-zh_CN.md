# qsctl

[![Build Status](https://travis-ci.org/qingstor/qsctl.svg?branch=master)](https://travis-ci.org/qingstor/qsctl)
[![GoDoc](https://godoc.org/github.com/qingstor/qsctl?status.svg)](https://godoc.org/github.com/qingstor/qsctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/qingstor/qsctl)](https://goreportcard.com/report/github.com/qingstor/qsctl)
[![codecov](https://codecov.io/gh/qingstor/qsctl/branch/master/graph/badge.svg)](https://codecov.io/gh/qingstor/qsctl)
[![localized](https://badges.crowdin.net/qsctl/localized.svg)](https://crowdin.com/project/qsctl)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/qingstor/qsctl/blob/master/LICENSE)
[![Join the Chat](https://img.shields.io/badge/chat-online-blue?style=flat&logo=mattermost)](https://chat.qingstor.dev/signup_user_complete/?id=1gjyqjsfo7dq7yfgomsjjg7h7o)

qsctl 是 QingStor 对象存储的高级命令行工具，它提供了强大的类 Unix 命令让你管理 QingStor 资源就像是在操作本地资源一般容易。

## 安装

### 二进制

访问 <https://github.com/qingstor/qsctl/releases> 以下载最新的文件。

## 快速开始

### 使用开始向导

### 手动配置

配置文件默认位于 `~/.qingstor/config.yaml`，也可以通过参数来手动指定 `-c /path/to/config`。

```yaml
access_key_id: 'ACCESS_KEY_ID_EXAMPLE'
secret_access_key: 'SECRET_ACCESS_KEY_EXAMPLE'
```

您也可以设置其他的配置项，如 `host`, `port` 等等，只需要在配置文件的下方按如下格式添加即可:

```yaml
host: 'qingstor.com'
port: 443
protocol: 'https'
connection_retries: 3
zone: 'zone_name'

# 默认情况下，对象键中的多余斜杠会被从路径中移除，
# 例如 `/a//b` 会变成 `a/b`，这个选项用于禁用该行为。
# disable_uri_cleaning: true
```

~~您也可以在没有配置文件的情况下执行 qsctl 指令，它首先会自动在特定的目录下寻找配置文件，
如果这些路径中都没有配置文件的话，qsctl 会启动一个交互式的配置程序来帮助您进行配置，
您只需要根据提示输入/选择配置内容即可。结束后会利用您输入的信息，在系统中生成配置文件
`{主目录}/.qingstor/config.yaml`。
(注意: 具体的配置文件的位置会根据您的系统而有所不同，通常来说，在类 Unix 操作系统下，
配置文件会生成在 `~/.qingstor/config.yaml`，而在 Windows 操作系统下，配置文件
会生成在 `C:\User\{用户名}\.qingstor\config.yaml`)~~
```
交互式配置从 v2.2.0 版本中被移除了，之后可能会以单独命令的方式添加回来。
```

从 v2.2.0 版本起，为了使您更好的在脚本中使用 `qsctl`，我们将所有的交互式操作移动到 `qsctl shell` 中了，
这使得脚本中执行的命令不会因为交互操作而中断。
同样地，我们也在命令执行时移除了进度条的渲染，并将其添加到了 `qsctl shell` 中了。
总之，现在您直接执行命令的话，不会有任何的交互行为。所有交互相关的行为都被移动到 `qsctl shell` 中了。


## qsctl shell

在 qsctl v2.2.0 版本中，我们新加入了交互式的 shell 界面，包含更多的引导和提示信息，推荐新用户使用。
您可以执行 qsctl shell 命令进入命令行界面，根据提示进行操作即可。
在命令行中，我们新增了对历史命令和自动补全的支持。其中：
- 在行开头可以自动提示补全可用命令；
- 输入 `qs://` 可以自动提示补全用户的 bucket 信息；
- 空格后可以自动提示补全本地文件信息；
- 输入 `-` 可以自动提示补全当前命令可用 flag 信息。


## 可用的命令

qsctl 支持如下命令

- `cat`: 获取指定对象的文件内容并输出到标准输出
- `cp`: 在本地和 QingStor 对象存储 bucket 之间复制数据
- `ls`: 列出 buckets 或者对象
- `mb`: 创建一个新的 bucket
- `mv`: 在本地和 QingStor 对象存储 bucket 之间移动数据
- `presign`: 获取指定对象的可公开访问链接
- `rb`: 删除一个 bucket
- `rm`: 删除一个或多个对象
- `stat`: 获取指定对象的信息
- `shell`: 启动一个 qsctl 的交互式 shell
- `sync`: 在本地和 QingStor 对象存储 bucket 之间进行同步
- `tee`: 从标准输入读取内容并上传为指定对象 (注意: qsctl 不会像 Linux 系统那样将内容 tee 操作到标准输出。)

详细的使用方法和样例可以通过 `qsctl --help` 或者 `qsctl <command> --help` 来获取。

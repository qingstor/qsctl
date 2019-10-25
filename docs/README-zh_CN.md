# qsctl

[![Build Status](https://travis-ci.org/yunify/qsctl.svg?branch=master)](https://travis-ci.org/yunify/qsctl)
[![GoDoc](https://godoc.org/github.com/yunify/qsctl?status.svg)](https://godoc.org/github.com/yunify/qsctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/yunify/qsctl)](https://goreportcard.com/report/github.com/yunify/qsctl)
[![codecov](https://codecov.io/gh/yunify/qsctl/branch/master/graph/badge.svg)](https://codecov.io/gh/yunify/qsctl)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/yunify/qsctl/blob/master/LICENSE)

qsctl 是 QingStor 对象存储的高级命令行工具，它提供了强大的类 Unix 命令让你管理 QingStor 资源就像是在操作本地资源一般容易。

## 安装

### 二进制

访问 <https://github.com/yunify/qsctl/releases> 以下载最新的文件。

## 快速开始

### 使用开始向导

### 手动配置

配置文件默认位于 `~/.qingstor/config.yaml`，也可以通过参数来手动 `-c /path/to/config`。

```yaml
access_key_id: 'ACCESS_KEY_ID_EXAMPLE'
secret_access_key: 'SECRET_ACCESS_KEY_EXAMPLE'
```

You can also config other option like `host` , `port` and so on, just
add lines below into configuration file, for example

```yaml
host: 'qingstor.com'
port: 443
protocol: 'https'
connection_retries: 3
# Valid levels are 'debug', 'info', 'warn', 'error', and 'fatal'.
log_level: 'debug'
zone: 'zone_name'
```

## 可用的命令

qsctl 支持如下命令

- `cat`: 获取指定对象的文件内容并输出到标准输出
- `cp`: 在本地和 QingStor 对象存储 bucket 之间复制数据
- `ls`: 列出 buckets 或者对象
- `mb`: 创建一个新的 bucket
- `presign`: 获取指定对象的可公开访问链接
- `rb`: 删除一个 bucket
- `rm`: 删除一个或多个对象
- `stat`: 获取指定对象的信息
- `sync`: 在本地和 QingStor 对象存储 bucket 之间进行同步
- `tee`: 从标准输入读取内容并上传为指定对象

详细的使用方法和样例可以通过 `qsctl help` 或者 `qsctl <command> help` 来获取。

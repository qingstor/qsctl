# -*- coding: utf-8 -*-
# =========================================================================
# Copyright (C) 2016 Yunify, Inc.
# -------------------------------------------------------------------------
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this work except in compliance with the License.
# You may obtain a copy of the License in the LICENSE file, or at:
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# =========================================================================

from __future__ import unicode_literals

import os
import sys
import argparse

from qingstor.sdk.config import Config
from qingstor.sdk.service.qingstor import QingStor

from ..constants import HTTP_OK, HTTP_OK_NO_CONTENT
from ..utils import (
    load_conf, uni_print, to_unix_path, is_pattern_match, validate_bucket_name,
    UploadIdRecorder
)


class BaseCommand(object):
    command = ""
    usage = ""
    description = ""

    client = None
    bucket_map = {}
    recorder = None
    options = None

    @classmethod
    def add_common_arguments(cls, parser):
        parser.add_argument(
            "-c",
            "--config",
            dest="config",
            action="store",
            type=str,
            default="~/.qingstor/config.yaml",
            help="Configuration file"
        )

    @classmethod
    def add_extra_arguments(cls, parser):
        pass

    @classmethod
    def add_transfer_arguments(cls, parser):
        pass

    @classmethod
    def get_argument_parser(cls):
        parser = argparse.ArgumentParser(
            prog="qsctl %s" % cls.command,
            usage=cls.usage,
            description=cls.description
        )
        cls.add_common_arguments(parser)
        cls.add_extra_arguments(parser)
        cls.add_transfer_arguments(parser)
        return parser

    @classmethod
    def get_client(cls, conf):
        config = Config().load_config_from_data(conf)
        return QingStor(config)

    @classmethod
    def get_zone(cls, bucket):
        url = "{protocol}://{bucket}.{host}:{port}".format(
            protocol=cls.client.config.protocol,
            host=cls.client.config.host,
            bucket=bucket,
            port=cls.client.config.port,
        )
        # cls.client.client is a Request Session
        resp = cls.client.client.head(url)
        if "Location" in resp.headers:
            # Location: http://test-bucket.zone.qingstor.com/
            zone = resp.headers["Location"].split(".")[1]
            return zone
        else:
            return ""

    @classmethod
    def send_request(cls):
        return None

    @classmethod
    def validate_bucket(cls, bucket):
        if not cls.bucket_map.get(bucket):
            cls.bucket_map[bucket] = cls.get_zone(bucket)
            if cls.bucket_map[bucket] == "":
                uni_print("Error: Please check if bucket <%s> exists" % bucket)
                sys.exit(-1)
            current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
            resp = current_bucket.head()
            if resp.status_code != HTTP_OK:
                uni_print(
                    "Error: Please check if you have enough"
                    " permission to access bucket <%s>." % bucket
                )
                sys.exit(-1)

    @classmethod
    def validate_local_path(cls, path):
        dirname = os.path.dirname(path)
        if dirname != "":
            if os.path.isfile(dirname):
                uni_print(
                    "Error: File with the same name '%s' already exists" %
                    dirname
                )
                sys.exit(-1)
            elif not os.path.isdir(dirname):
                try:
                    os.makedirs(dirname)
                    uni_print("Directory '%s' created" % dirname)
                except OSError as e:
                    uni_print(
                        "Error: Failed to create directory '%s': %s" %
                        (dirname, e)
                    )
                    sys.exit(-1)

    @classmethod
    def validate_qs_path(cls, qs_path):
        qs_path = to_unix_path(qs_path)
        if qs_path.startswith("qs://"):
            qs_path = qs_path[5:]
        qs_path_split = qs_path.split("/", 1)
        if len(qs_path_split) == 1:
            bucket, prefix = qs_path_split[0], ""
        elif len(qs_path_split) == 2:
            bucket, prefix = qs_path_split[0], qs_path_split[1]
        if not validate_bucket_name(bucket):
            uni_print("Error: Invalid Bucket name")
            sys.exit(-1)
        if cls.command not in ("mb", "rb"):
            cls.validate_bucket(bucket)
        return bucket, prefix

    @classmethod
    def key_exists(cls, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.head_object(key)
        return resp.status_code == HTTP_OK

    @classmethod
    def remove_key(cls, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.head_object(key)
        if resp.status_code == HTTP_OK:
            if resp.headers["Content-Type"] == "application/x-directory":
                statement = "Directory should be deleted with -r"
                uni_print(statement)
            else:
                resp = current_bucket.delete_object(key)
                if resp.status_code == HTTP_OK_NO_CONTENT:
                    statement = "Key <%s> deleted" % key
                    uni_print(statement)
                else:
                    uni_print(resp.content)
        else:
            statement = "Key <%s> does not exist" % key
            uni_print(statement)

    @classmethod
    def confirm_key_remove(cls, key_name):
        if cls.command == "rb":
            return True
        else:
            return is_pattern_match(
                key_name, cls.options.exclude, cls.options.include
            )

    @classmethod
    def remove_multiple_keys(cls, bucket, prefix=""):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        marker = ""
        while True:
            keys, marker, _ = cls.list_multiple_keys(
                bucket, marker=marker, prefix=prefix, limit="1000"
            )
            keys_to_remove = [i["key"] for i in keys]
            for key in keys_to_remove:
                if not cls.confirm_key_remove(key[len(prefix):]):
                    keys_to_remove.remove(key)
            keys_to_remove = [{"key": key} for key in keys_to_remove]
            resp = current_bucket.delete_multiple_objects(
                objects=keys_to_remove
            )
            if resp.status_code == HTTP_OK:
                keys_removed = [i["key"] for i in resp["deleted"]]
                for key in keys_removed:
                    statement = "Key <%s> deleted" % key
                    uni_print(statement)
                keys_error = resp["errors"]
                for key in keys_error:
                    statement = "Key <%s> deleted failed for <%s> " % (
                        key["key"], key["message"]
                    )
                    uni_print(statement)
            else:
                uni_print(resp.content)
            if marker == "":
                break

    @classmethod
    def list_multiple_keys(
            cls, bucket, prefix="", delimiter="", marker="", limit="200"
    ):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.list_objects(
            marker=marker, prefix=prefix, delimiter=delimiter, limit=limit
        )
        keys = resp["keys"]
        dirs = resp["common_prefixes"]
        next_marker = resp["next_marker"]
        return keys, next_marker, dirs

    @classmethod
    def main(cls, args):

        parser = cls.get_argument_parser()
        cls.options = parser.parse_args(args)

        # Load config file
        config_path = ["~/.qingstor/config.yaml", "~/.qingcloud/config.yaml"]

        # IF has options.config, insert it
        config_path.insert(0, cls.options.config)

        for path in config_path:
            conf = load_conf(path)
            if conf is not None:
                # Get client of qingstor
                cls.client = cls.get_client(conf)
                break

        cls._init_recorder()

        if cls.client is None:
            sys.exit(-1)

        # Send request
        return cls.send_request()

    @classmethod
    def _init_recorder(cls):
        # Init UploadIdRecorder
        qsctl_dir = os.path.expanduser("~/.qingstor/qsctl")
        if not os.path.exists(qsctl_dir):
            os.makedirs(qsctl_dir)
        record_path = os.path.join(qsctl_dir, "record")
        cls.recorder = UploadIdRecorder(record_path)

    @classmethod
    def _handle_sigint(cls, signature, frame):
        # Handler function for signal.SIGINT
        if cls.recorder:
            cls.recorder.close()
        sys.exit(0)

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

import os
import sys
import argparse

from qingstor.sdk.config import Config
from qingstor.sdk.service.qingstor import QingStor

from ..constants import HTTP_OK, HTTP_OK_NO_CONTENT
from ..utils import (
    load_conf,
    uni_print,
    to_unix_path,
    is_pattern_match,
    validate_bucket_name,
)


class BaseCommand(object):
    command = ""
    usage = ""
    description = ""

    client = None
    bucket_map = {}

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
            prog='qsctl %s' % cls.command,
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
    def get_buckets(cls):
        resp = cls.client.list_buckets()
        if resp.status_code != HTTP_OK:
            print(
                "Error: Please check your configuration and you have "
                "enough permission to access qingstor service."
            )
            sys.exit()
        return resp["buckets"]

    @classmethod
    def send_request(cls, options):
        return None

    @classmethod
    def validate_bucket(cls, bucket):
        if not cls.bucket_map.get(bucket):
            cls.bucket_map[bucket] = ""
            buckets = cls.get_buckets()
            for i in buckets:
                if bucket == i["name"]:
                    cls.bucket_map[bucket] = i["location"]
                    break
            if cls.bucket_map[bucket] == "":
                print("Error: Please check if bucket <%s> exists" % bucket)
                sys.exit(-1)
            current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
            resp = current_bucket.head()
            if resp.status_code != HTTP_OK:
                print(
                    "Error: Please check if you have enough"
                    " permission to access bucket <%s>." % bucket
                )
                sys.exit(-1)

    @classmethod
    def validate_local_path(cls, path):
        dirname = os.path.dirname(path)
        if dirname != "":
            if os.path.isfile(dirname):
                print(
                    "Error: File with the same name '%s' already exists" %
                    dirname
                )
                sys.exit(-1)
            elif not os.path.isdir(dirname):
                try:
                    os.makedirs(dirname)
                    print("Directory '%s' created" % dirname)
                except OSError as e:
                    print(
                        "Error: Failed to create directory '%s': %s" %
                        (dirname, e)
                    )
                    sys.exit(-1)

    @classmethod
    def validate_qs_path(cls, qs_path):
        qs_path = to_unix_path(qs_path)
        if qs_path.startswith("qs://"):
            qs_path = qs_path[5:]
        qs_path_split = qs_path.split('/', 1)
        if len(qs_path_split) == 1:
            bucket, prefix = qs_path_split[0], ""
        elif len(qs_path_split) == 2:
            bucket, prefix = qs_path_split[0], qs_path_split[1]
        if not validate_bucket_name(bucket):
            print("Error: Invalid Bucket name")
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
        resp = current_bucket.delete_object(key)
        if resp.status_code != HTTP_OK_NO_CONTENT:
            print(resp.content)
        else:
            statement = "Key <%s> deleted" % key
            uni_print(statement)

    @classmethod
    def confirm_key_remove(cls, key_name, options):
        if cls.command == "rb":
            return True
        else:
            return is_pattern_match(key_name, options.exclude, options.include)

    @classmethod
    def remove_multiple_keys(cls, bucket, prefix="", options=None):
        cls.validate_bucket(bucket)
        marker = ""
        while True:
            keys, marker, _ = cls.list_multiple_keys(
                bucket, marker=marker, prefix=prefix
            )
            for item in keys:
                key = item["key"] if sys.version > "3" else item["key"].encode(
                    'utf8'
                )
                if cls.confirm_key_remove(key[len(prefix):], options):
                    cls.remove_key(bucket, key)
            if marker == "":
                break

    @classmethod
    def list_multiple_keys(
            cls, bucket, prefix="", delimiter="", marker="", limit=200
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
        options = parser.parse_args(args)

        # Load config file
        conf = load_conf(options.config)

        if conf is None:
            sys.exit(-1)

        # Get client of qingstor
        cls.client = cls.get_client(conf)

        # Send request
        return cls.send_request(options)

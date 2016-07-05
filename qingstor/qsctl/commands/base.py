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

from qingcloud.qingstor.connection import QSConnection

from ..constants import ENDPOINT, HTTP_OK, HTTP_OK_NO_CONTENT
from ..utils import (
    load_conf,
    is_pattern_match,
    to_unix_path,
    uni_print,
    json_loads,
)

class BaseCommand(object):

    command = ""
    usage = ""
    description = ""

    conn = None

    @classmethod
    def add_common_arguments(cls, parser):
        parser.add_argument(
            "-c",
            "--config",
            dest="config",
            action="store",
            type=str,
            default="~/.qingcloud/config.yaml",
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
    def get_connection(cls, conf, options):
        if cls.command in ("mb"):
            host = ENDPOINT
        else:
            host = "%s.%s" % (conf["zone"], ENDPOINT)

        return QSConnection(
            qy_access_key_id=conf["qy_access_key_id"],
            qy_secret_access_key=conf["qy_secret_access_key"],
            host=host
        )

    @classmethod
    def send_request(cls, options):
        return None

    @classmethod
    def validate_bucket(cls, bucket):
        resp = cls.conn.make_request("HEAD", bucket)
        status = resp.status
        resp.close()
        if status != HTTP_OK:
            print("Error: Please check if bucket <%s> exists and you have " \
                "enough permission to access it." % bucket)
            sys.exit(-1)

    @classmethod
    def validate_local_path(cls, path):
        dirname = os.path.dirname(path)
        if dirname != "":
            if os.path.isfile(dirname):
                print("Error: File with the same name '%s' already exists" \
                    % dirname)
                sys.exit(-1)
            elif not os.path.isdir(dirname):
                try:
                    os.makedirs(dirname)
                    print("Directory '%s' created" % dirname)
                except OSError as e:
                    print("Error: Failed to create directory '%s': %s" % (dirname, e))
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
        if cls.command not in ("mb", "rb"):
            cls.validate_bucket(bucket)
        return bucket, prefix

    @classmethod
    def key_exists(cls, bucket, key):
        resp = cls.conn.make_request("HEAD", bucket, key)
        status = resp.status
        resp.close()
        return status == HTTP_OK

    @classmethod
    def remove_key(cls, bucket, key):
        resp = cls.conn.make_request("DELETE", bucket, key)
        if resp.status != HTTP_OK_NO_CONTENT:
            print(resp.read())
        else:
            statement = "Key <%s> deleted" % key
            uni_print(statement)
        resp.close()

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
            keys, marker, _ = cls.list_multiple_keys(bucket, marker, prefix)
            for item in keys:
                key = item["key"] if sys.version > "3" else item["key"].encode('utf8')
                if cls.confirm_key_remove(key[len(prefix):], options):
                    cls.remove_key(bucket, key)
            if marker == "":
                break

    @classmethod
    def list_multiple_keys(cls, bucket, marker, prefix="", params=None):
        cls.validate_bucket(bucket)
        if params == None:
            params = {}
        if prefix != "":
            params["prefix"] = prefix
        if marker != "" and cls.key_exists(bucket, marker):
            params["marker"] = marker
        resp = cls.conn.make_request("GET", bucket, params=params)
        body = json_loads(resp.read())
        keys = body["keys"]
        dirs = body["common_prefixes"]
        marker = body["marker"] if sys.version > "3" else body["marker"].encode('utf8')
        resp.close()
        return keys, marker, dirs

    @classmethod
    def main(cls, args):

        parser = cls.get_argument_parser()
        options = parser.parse_args(args)

        # Load config file
        conf = load_conf(options.config)
        if conf is None:
            sys.exit(-1)

        # Get connection to the server
        cls.conn = cls.get_connection(conf, options)

        # Send request
        return cls.send_request(options)

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
import time

from .transfer import TransferCommand

from ..constants import HTTP_OK
from ..utils import is_pattern_match, to_unix_path, join_local_path, uni_print


class SyncCommand(TransferCommand):

    command = "sync"
    usage = (
        "%(prog)s <source-path> <dest-path> [-c <conf_file> --delete "
        "--exclude <pattern value> --include <pattern value> --rate-limit <pattern value>]"
    )

    @classmethod
    def add_transfer_arguments(cls, parser):
        parser.add_argument(
            "--delete",
            action="store_true",
            dest="delete",
            help=(
                "Any files or keys existing under the destination directory "
                "but not existing in the source directory will be deleted if "
                "--delete option is specified."
            )
        )
        return parser

    @classmethod
    def cleanup(cls, transfer_flow, bucket, prefix):
        if cls.options.delete == True:
            if transfer_flow == "LOCAL_TO_QS":
                cls.clean_keys(bucket, prefix)
            elif transfer_flow == "QS_TO_LOCAL":
                cls.clean_files(bucket, prefix)

    @classmethod
    def clean_files(cls, bucket, prefix):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        for rt, dirs, files in os.walk(cls.options.dest_path):
            for f in files:
                local_path = os.path.join(rt, f)
                key_path = os.path.relpath(local_path, cls.options.dest_path)
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                resp = current_bucket.head_object(key)
                if (resp.status_code != HTTP_OK) or (
                        not is_pattern_match(
                            key_path, cls.options.exclude, cls.options.include
                        )
                ):
                    os.remove(local_path)
                    uni_print("File '%s' deleted" % local_path)

        for rt, dirs, files in os.walk(cls.options.dest_path):
            for d in dirs:
                local_path = os.path.join(rt, d)
                key_path = os.path.relpath(
                    local_path, cls.options.dest_path
                ) + "/"
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                resp = current_bucket.head_object(key)
                if (resp.status_code != HTTP_OK) or (
                        not is_pattern_match(
                            key_path, cls.options.exclude, cls.options.include
                        )
                ):
                    if not os.listdir(local_path):
                        os.rmdir(local_path)
                        uni_print("Directory '%s' deleted" % local_path)

    @classmethod
    def clean_keys(cls, bucket, prefix):
        cls.remove_multiple_keys(bucket, prefix)

    @classmethod
    def confirm_key_upload(cls, local_path, bucket, key):
        if cls.key_exists(bucket, key):
            time_key_modified = cls.get_time_key_modified(bucket, key)
            return cls.is_local_file_modified(local_path, time_key_modified)
        else:
            return True

    @classmethod
    def get_time_key_modified(cls, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.head_object(key)
        if resp.status_code == HTTP_OK:
            time_str_key = resp.headers["Last-Modified"]
            return time.mktime(
                time.strptime(time_str_key, "%a, %d %b %Y %X GMT")
            )
        else:
            statement = "Error: Failed to head key <%s>" % key
            uni_print(statement)
            sys.exit(-1)

    @classmethod
    def is_local_file_modified(cls, local_path, time_key_modified):
        time_stamp_file = os.stat(local_path).st_mtime
        time_file_modified = time.mktime(time.gmtime(time_stamp_file))
        return time_file_modified > time_key_modified

    @classmethod
    def confirm_key_download(cls, local_path, time_key_modified=None):
        if os.path.isfile(local_path):
            time_file_modified = os.stat(local_path).st_mtime
            return time_key_modified > time_file_modified
        else:
            return True

    @classmethod
    def confirm_key_remove(cls, key):
        file_path = join_local_path(cls.options.source_path, key)
        return (not os.path.exists(file_path)) or (
            not is_pattern_match(key, cls.options.exclude, cls.options.include)
        )

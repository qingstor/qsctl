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
import errno

from .base import BaseCommand

from ..constants import (
    PART_SIZE,
    HTTP_OK,
    HTTP_OK_CREATED,
    HTTP_OK_PARTIAL_CONTENT,
)

from ..utils import (
    confirm_by_user,
    is_pattern_match,
    to_unix_path,
    join_local_path,
    uni_print,
    get_part_numbers,
    FileChunk,
    StdinFileChunk,
)


class TransferCommand(BaseCommand):

    command = ""
    usage = ""

    @classmethod
    def add_extra_arguments(cls, parser):
        parser.add_argument("source_path", help="Original path")

        parser.add_argument("dest_path", help="Destination path")

        parser.add_argument(
            "--exclude",
            type=str,
            help="Exclude all files or keys that match the specified pattern."
        )

        parser.add_argument(
            "--include",
            type=str,
            help="Do not exclude files or keys that match the specified pattern"
        )

    @classmethod
    def add_transfer_arguments(cls, parser):
        parser.add_argument(
            "-r",
            "--recursive",
            action="store_true",
            dest="recursive",
            help="Recursively transfer keys"
        )

        parser.add_argument(
            "-f",
            "--force",
            action="store_true",
            dest="force",
            help="Forcely overwrite existing key or file without asking"
        )
        return parser

    @classmethod
    def get_transfer_flow(cls, options):
        source, dest = options.source_path, options.dest_path
        if source.startswith("qs://") and not (dest.startswith("qs://")):
            return "QS_TO_LOCAL"
        elif dest.startswith("qs://") and not (source.startswith("qs://")):
            return "LOCAL_TO_QS"
        else:
            print(
                "Error: please give correct local path and qs-path. "
                "The qs_path must start with 'qs://'."
            )
            sys.exit(-1)

    @classmethod
    def confirm_key_upload(cls, options, local_path, bucket, key):
        if options.force or not cls.key_exists(bucket, key):
            return True
        else:
            notice = "Key <%s> already existed in bucket <%s>.\n" % (
                key, bucket
            )
            notice += "Overwrite key <%s>? (y or n) " % key
            return confirm_by_user(notice)

    @classmethod
    def confirm_key_download(cls, options, local_path, time_key_modified=None):
        if options.force or not os.path.isfile(local_path):
            return True
        else:
            notice = "File '%s' already existed.\n" % local_path
            notice += "Overwrite file '%s'? (y or n) " % local_path
            return confirm_by_user(notice)

    @classmethod
    def upload_files(cls, options):
        if not os.path.isdir(options.source_path):
            print("Error: No such directory: %s" % options.source_path)
            sys.exit(-1)

        bucket, prefix = cls.validate_qs_path(options.dest_path)
        if prefix != "" and (not prefix.endswith("/")):
            prefix = prefix + "/"

        for rt, dirs, files in os.walk(options.source_path):
            for d in dirs:
                local_path = os.path.join(rt, d)
                key_path = os.path.relpath(
                    local_path, options.source_path
                ) + "/"
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                if (is_pattern_match(key_path, options.exclude, options.include)
                        and cls.confirm_key_upload(options, local_path, bucket,
                                                   key)):
                    cls.put_directory(bucket, key)

            for f in files:
                local_path = os.path.join(rt, f)
                key_path = os.path.relpath(local_path, options.source_path)
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                if (is_pattern_match(key_path, options.exclude, options.include)
                        and cls.confirm_key_upload(options, local_path, bucket,
                                                   key)):
                    cls.send_local_file(local_path, bucket, key)

        cls.cleanup("LOCAL_TO_QS", options, bucket, prefix)

    @classmethod
    def upload_file(cls, options):
        bucket, prefix = cls.validate_qs_path(options.dest_path)
        if prefix.endswith("/") or prefix == "":
            key = prefix + os.path.basename(options.source_path)
        else:
            key = prefix
        if os.path.isfile(options.source_path):
            if cls.confirm_key_upload(
                    options, options.source_path, bucket, key
            ):
                cls.send_local_file(options.source_path, bucket, key)
        elif options.source_path == '-':
            cls.send_data_from_stdin(bucket, key)
        else:
            print("Error: No such file: %s" % options.source_path)
            sys.exit(-1)

    @classmethod
    def download_files(cls, options):
        bucket, prefix = cls.validate_qs_path(options.source_path)
        if prefix != "" and (not prefix.endswith("/")):
            prefix = prefix + "/"
        marker = ""
        while True:
            keys, marker, _ = cls.list_multiple_keys(
                bucket, marker=marker, prefix=prefix
            )
            for item in keys:
                key = item["key"] if sys.version > "3" else item["key"].encode(
                    'utf8'
                )
                key_name = key[len(prefix):]
                local_path = join_local_path(options.dest_path, key_name)
                if (is_pattern_match(key_name, options.exclude, options.include)
                        and cls.confirm_key_download(options, local_path,
                                                     item["modified"])):
                    cls.write_local_file(local_path, bucket, key)
            if marker == "":
                break

        cls.cleanup("QS_TO_LOCAL", options, bucket, prefix)

    @classmethod
    def download_file(cls, options):
        bucket, key = cls.validate_qs_path(options.source_path)
        if key == "":
            print(
                "Error: Please give correct and complete key qs-path, such "
                "as 'qs://yourbucket/key'."
            )
            sys.exit(-1)
        if os.path.isdir(options.dest_path):
            local_path = join_local_path(options.dest_path, key)
        else:
            local_path = options.dest_path
        if cls.confirm_key_download(options, local_path):
            cls.write_local_file(local_path, bucket, key)

    @classmethod
    def write_local_file(cls, local_path, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.get_object(key)
        if resp.status_code in (HTTP_OK, HTTP_OK_PARTIAL_CONTENT):
            cls.validate_local_path(local_path)
            if key[-1] != "/":
                with open(local_path, "wb") as f:
                    for chunk in resp.iter_content():
                        f.write(chunk)
                print("File '%s' written" % local_path)
            if cls.command == "mv":
                cls.remove_key(bucket, key)
        else:
            print(resp.content)
            sys.exit(-1)

    @classmethod
    def put_directory(cls, bucket, key):
        content_type = "application/x-directory"
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.put_object(key, content_type=content_type)
        if resp.status_code == HTTP_OK_CREATED:
            statement = "Directory <%s> created in bucket <%s>" % (key, bucket)
            uni_print(statement)
        else:
            print(resp.content)

    @classmethod
    def send_local_file(cls, local_path, bucket, key):
        try:
            if os.path.getsize(local_path) > PART_SIZE:
                cls.multipart_upload_file(local_path, bucket, key)
            else:
                cls.send_file(local_path, bucket, key)
        except OSError as e:
            if e.errno == errno.ENOENT:
                print(
                    "WARN: file %s not found, perhaps it's removed during "
                    "qsctl operation" % local_path
                )

    @classmethod
    def upload_multipart_from_stdin(cls, upload_id, bucket, key):
        global upload_failed
        global part_numbers
        part_numbers = []
        next_part_number = 0
        while True:
            if upload_failed == True:
                break

            data = StdinFileChunk(PART_SIZE)
            done = len(data)
            part_numbers.append(next_part_number)
            cls.upload_part(upload_id, next_part_number, data, bucket, key)
            next_part_number = next_part_number + 1

            if done != PART_SIZE:
                break

    @classmethod
    def send_data_from_stdin(cls, bucket, key):
        global upload_failed
        upload_failed = False
        upload_id = cls.init_multipart(bucket, key)

        if upload_id == "":
            statement = "Error: key <%s> already exists" % key
            uni_print(statement)
        else:
            cls.upload_multipart_from_stdin(upload_id, bucket, key)
            if not upload_failed:
                cls.complete_multipart('-', upload_id, bucket, key)
            else:
                print("Error: Failed to upload file '%s'" % '-')

    @classmethod
    def send_file(cls, local_path, bucket, key):
        with open(local_path, "rb") as data:
            cls.validate_bucket(bucket)
            current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
            resp = current_bucket.put_object(key, body=data)
        if resp.status_code == HTTP_OK_CREATED:
            statement = "Key <%s> created in bucket <%s>" % (key, bucket)
            uni_print(statement)
            if cls.command == "mv":
                os.remove(local_path)
        else:
            print(resp.content)

    @classmethod
    def multipart_upload_file(cls, local_path, bucket, key):
        global upload_failed
        upload_failed = False
        upload_id = cls.init_multipart(bucket, key)
        if upload_id == "":
            statement = "Error: key <%s> already exists" % key
            uni_print(statement)
        else:
            cls.upload_multipart(local_path, upload_id, bucket, key)
            if not upload_failed:
                cls.complete_multipart(local_path, upload_id, bucket, key)
            else:
                print("Error: Failed to upload file '%s'" % local_path)

    @classmethod
    def init_multipart(cls, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.initiate_multipart_upload(key)
        if resp.status_code != HTTP_OK:
            upload_id = ""
        else:
            upload_id = resp["upload_id"]
        return upload_id

    @classmethod
    def upload_multipart(cls, filepath, upload_id, bucket, key):
        global upload_failed
        global part_numbers
        part_numbers = get_part_numbers(filepath)
        for part_number in part_numbers:
            if upload_failed == True:
                break

            data = FileChunk(filepath, part_number)
            cls.upload_part(upload_id, part_number, data, bucket, key)

    @classmethod
    def upload_part(cls, upload_id, part_number, data, bucket, key):
        '''upload one part of large file.
        '''
        global upload_failed

        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.upload_multipart(
            key, upload_id=upload_id, part_number=part_number, body=data
        )
        if resp.status_code != HTTP_OK_CREATED:
            print(resp.content)
            upload_failed = True
        data.close()

    @classmethod
    def complete_multipart(cls, filepath, upload_id, bucket, key):
        global part_numbers
        parts = []
        for part_number in part_numbers:
            parts.append({"part_number": part_number})
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.complete_multipart_upload(
            key, upload_id=upload_id, object_parts=parts
        )
        if resp.status_code != HTTP_OK_CREATED:
            print(resp.content)
        else:
            statement = "Key <%s> created in bucket <%s>" % (key, bucket)
            uni_print(statement)
            if cls.command == "mv":
                os.remove(filepath)

    @classmethod
    def cleanup(cls, transfer_flow, options, bucket, prefix):
        pass

    @classmethod
    def send_request(cls, options):
        transfer_flow = cls.get_transfer_flow(options)
        if transfer_flow == "LOCAL_TO_QS":
            if (cls.command == "sync") or (options.recursive is True):
                cls.upload_files(options)
            else:
                cls.upload_file(options)
        elif transfer_flow == "QS_TO_LOCAL":
            if (cls.command == "sync") or (options.recursive is True):
                cls.download_files(options)
            else:
                cls.download_file(options)

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
import errno
import signal

from tqdm import tqdm
from .base import BaseCommand

from ..constants import (
    PART_SIZE,
    BAR_FORMAT,
    USE_ASCII,
    HTTP_OK,
    HTTP_OK_CREATED,
    HTTP_BAD_REQUEST,
    HTTP_OK_PARTIAL_CONTENT,
    TEMPORARY_FILE_SUFFIX,
)

from ..utils import (
    confirm_by_user, is_pattern_match, to_unix_path, join_local_path, uni_print,
    FileChunk, wrapper_stream, convert_to_bytes, TokenPail
)


class TransferCommand(BaseCommand):
    command = ""
    usage = ""

    tokens = None

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

        parser.add_argument(
            "--no-progress",
            action="store_true",
            default=False,
            dest="no_progress",
            help="Close progress bar display"
        )

        parser.add_argument(
            "--rate-limit",
            type=str,
            help="add rate limit for a second, eg: --rate-limit 500K"
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
    def get_transfer_flow(cls):
        source, dest = cls.options.source_path, cls.options.dest_path
        if source.startswith("qs://") and not (dest.startswith("qs://")):
            return "QS_TO_LOCAL"
        elif dest.startswith("qs://") and not (source.startswith("qs://")):
            return "LOCAL_TO_QS"
        else:
            uni_print(
                "Error: please give correct local path and qs-path. "
                "The qs_path must start with 'qs://'."
            )
            sys.exit(-1)

    @classmethod
    def confirm_key_upload(cls, local_path, bucket, key):
        if key.endswith(TEMPORARY_FILE_SUFFIX):
            # Skip temporary file created in downloading process
            return False

        if cls.options.force or not cls.key_exists(bucket, key):
            return True
        else:
            notice = "Key <%s> already existed in bucket <%s>.\n" % (
                key, bucket
            )
            notice += "Overwrite key <%s>? (y or n) " % key
            return confirm_by_user(notice)

    @classmethod
    def confirm_key_download(cls, local_path, time_key_modified=None):
        if cls.options.force or not os.path.isfile(local_path):
            return True
        else:
            notice = "File '%s' already existed.\n" % local_path
            notice += "Overwrite file '%s'? (y or n) " % local_path
            return confirm_by_user(notice)

    @classmethod
    def upload_files(cls):
        source_path = cls.options.source_path
        dest_path = cls.options.dest_path
        if not os.path.isdir(source_path):
            uni_print("Error: No such directory: %s" % source_path)
            sys.exit(-1)

        bucket, prefix = cls.validate_qs_path(dest_path)
        if prefix != "" and (not prefix.endswith("/")):
            prefix += "/"

        for rt, dirs, files in os.walk(source_path):
            for d in dirs:
                local_path = os.path.join(rt, d)
                key_path = os.path.relpath(local_path, source_path) + "/"
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                if (is_pattern_match(key_path, cls.options.exclude, cls.options.include)
                    and cls.confirm_key_upload(local_path, bucket,
                                               key)):
                    cls.put_directory(bucket, key)

            for f in files:
                local_path = os.path.join(rt, f)
                key_path = os.path.relpath(local_path, source_path)
                key_path = to_unix_path(key_path)
                key = prefix + key_path
                if (is_pattern_match(key_path, cls.options.exclude, cls.options.include)
                    and cls.confirm_key_upload(local_path, bucket,
                                               key)):
                    cls.send_local_file(local_path, bucket, key)

        cls.cleanup("LOCAL_TO_QS", bucket, prefix)

    @classmethod
    def upload_file(cls):
        source_path = cls.options.source_path
        dest_path = cls.options.dest_path
        bucket, prefix = cls.validate_qs_path(dest_path)
        if prefix.endswith("/") or prefix == "":
            key = prefix + os.path.basename(source_path)
        else:
            key = prefix
        if os.path.isfile(source_path):
            if cls.confirm_key_upload(source_path, bucket, key):
                cls.send_local_file(source_path, bucket, key)
        elif source_path == '-':
            cls.send_data_from_stdin(bucket, key)
        else:
            uni_print("Error: No such file: %s" % source_path)
            sys.exit(-1)

    @classmethod
    def download_files(cls):
        source_path = cls.options.source_path
        dest_path = cls.options.dest_path
        bucket, prefix = cls.validate_qs_path(source_path)
        if prefix != "" and (not prefix.endswith("/")):
            prefix += "/"
        marker = ""
        while True:
            keys, marker, _ = cls.list_multiple_keys(
                bucket, marker=marker, prefix=prefix
            )
            for item in keys:
                key = item["key"]
                key_name = key[len(prefix):]
                local_path = join_local_path(dest_path, key_name)
                is_match = is_pattern_match(
                    key_name, cls.options.exclude, cls.options.include
                )
                is_confirmed_key_download = cls.confirm_key_download(
                    local_path, item["modified"]
                )
                if local_path and is_match and is_confirmed_key_download:
                    cls.write_local_file(local_path, bucket, key)
            if marker == "":
                break

        cls.cleanup("QS_TO_LOCAL", bucket, prefix)

    @classmethod
    def download_file(cls):
        source_path = cls.options.source_path
        dest_path = cls.options.dest_path
        bucket, key = cls.validate_qs_path(source_path)
        if key == "":
            uni_print(
                "Error: Please give correct and complete key qs-path, such "
                "as 'qs://yourbucket/key'."
            )
            sys.exit(-1)
        if os.path.isdir(dest_path):
            local_path = join_local_path(dest_path, key)
        else:
            local_path = dest_path
        if local_path and cls.confirm_key_download(local_path):
            cls.write_local_file(local_path, bucket, key)

    @classmethod
    def write_local_file(cls, local_path, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])

        completed = 0
        temporary_path = local_path + TEMPORARY_FILE_SUFFIX
        if os.path.isfile(temporary_path):
            completed = os.path.getsize(temporary_path)

        if completed > 0:
            resp = current_bucket.get_object(key, range="bytes=%d-" % completed)
            uni_print("Resume downloading key <%s>" % key)
        else:
            resp = current_bucket.get_object(key)
        if resp.status_code in (HTTP_OK, HTTP_OK_PARTIAL_CONTENT):
            cls.validate_local_path(local_path)
            open_flag = "wb"
            if key[-1] != "/":
                content_length = int(resp.headers["Content-Length"])
                if completed > 0:
                    open_flag = "ab"

                with open(temporary_path, open_flag) as f:
                    uni_print(
                        "Key <%s> is downloading as File <%s>" %
                        (key, temporary_path)
                    )
                    if cls.options.no_progress:
                        pbar = None
                    else:
                        pbar = tqdm(
                            initial=completed,
                            total=completed + content_length,
                            unit="B",
                            unit_scale=True,
                            desc="Transferring",
                            bar_format=BAR_FORMAT,
                            leave=False,
                            ascii=USE_ASCII
                        )
                    cache = []

                    for chunk in resp.iter_content(1024):
                        # cls.tokens is not None , rate limit
                        if cls.get_tokens_obj():
                            while not cls.get_tokens_obj().consume(1024):
                                continue
                        if pbar:
                            pbar.update(1024)
                        cache.append(chunk)
                        # Write file while cache is over 32M
                        if len(cache) >= 32 * 1024:
                            f.write(b"".join(cache))
                            cache = []

                    if cache:
                        f.write(b"".join(cache))
                        del cache
                    if pbar is not None:
                        pbar.close()

                    os.rename(temporary_path, local_path)
                    uni_print(
                        "File <%s> written, rename to original file <%s>" %
                        (temporary_path, local_path),
                    )

            if cls.command == "mv":
                cls.remove_key(bucket, key)
        else:
            uni_print(resp.content)
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
            uni_print(resp.content)

    @classmethod
    def send_local_file(cls, local_path, bucket, key):
        try:
            if os.path.getsize(local_path) > PART_SIZE:
                cls.multipart_upload_file(local_path, bucket, key)
            else:
                cls.send_file(local_path, bucket, key)
        except OSError as e:
            if e.errno == errno.ENOENT:
                uni_print(
                    "WARN: file %s not found, perhaps it's removed during "
                    "qsctl operation" % local_path
                )

    @classmethod
    def send_data_from_stdin(cls, bucket, key):
        fc = FileChunk(sys.stdin)
        try:
            (_, data) = next(fc.iter())
        except StopIteration:
            return
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        uni_print("Stdin is uploading as Key <%s>" % (key))
        if cls.options.no_progress:
            pbar = None
        else:
            pbar = tqdm(
                total=fc.size,
                unit="B",
                unit_scale=True,
                desc="Transferring",
                bar_format=BAR_FORMAT,
                leave=False,
                ascii=USE_ASCII
            )
        resp = current_bucket.put_object(
            key, body=wrapper_stream(data, pbar, cls.get_tokens_obj())
        )
        if pbar:
            pbar.close()
        if resp.status_code == HTTP_OK_CREATED:
            statement = "Key <%s> created in bucket <%s>" % (key, bucket)
            uni_print(statement)
        else:
            uni_print(resp.content)

    @classmethod
    def send_file(cls, local_path, bucket, key):
        with open(local_path, "rb") as f:
            fc = FileChunk(f)
            try:
                (_, data) = next(fc.iter())
            except StopIteration:
                return
            cls.validate_bucket(bucket)
            current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
            uni_print("File <%s> is uploading as Key <%s>" % (local_path, key))
            if cls.options.no_progress:
                pbar = None
            else:
                pbar = tqdm(
                    total=fc.size,
                    unit="B",
                    unit_scale=True,
                    desc="Transferring",
                    bar_format=BAR_FORMAT,
                    leave=False,
                    ascii=USE_ASCII
                )
            resp = current_bucket.put_object(
                key, body=wrapper_stream(data, pbar, cls.get_tokens_obj())
            )
            if pbar:
                pbar.close()
            if resp.status_code == HTTP_OK_CREATED:
                statement = "Key <%s> created in bucket <%s>" % (key, bucket)
                uni_print(statement)
                if cls.command == "mv":
                    os.remove(local_path)
            else:
                uni_print(resp.content)

    @classmethod
    def multipart_upload_file(cls, local_path, bucket, key):
        local_path = os.path.join(os.getcwd(), local_path)
        resume_multipart, upload_id, cur_part_number = \
            cls.try_to_resume_multipart(local_path, bucket, key)
        if not resume_multipart:
            upload_id = cls.init_multipart(bucket, key)
            cls.recorder.put_record(local_path, bucket, key, upload_id)
        is_upload_success, cur_parts = cls.upload_multipart(
            local_path, upload_id, bucket, key, cur_part_number
        )
        if is_upload_success:
            cls.complete_multipart(
                local_path, upload_id, cur_parts, bucket, key
            )
        else:
            uni_print("Error: Failed to upload file <%s>" % local_path)

    @classmethod
    def try_to_resume_multipart(cls, local_path, bucket, key):
        upload_id = cls.recorder.get_record(local_path, bucket, key)
        if not upload_id:
            return False, "", 0
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.list_multipart(key, upload_id=upload_id)
        if resp.status_code == HTTP_BAD_REQUEST:
            # Previous upload has been aborted or completed.
            uni_print(
                "Warning: Previous upload has been aborted or completed. "
                "Can't resume uploading key <%s> via previous upload id <%s>" %
                (key, upload_id)
            )
            cls.recorder.remove_record(local_path, bucket, key)
            return False, "", 0
        elif resp.status_code != HTTP_OK:
            uni_print(
                "Failed to list multipart. Response code: %d. "
                "Response content: %s" % (resp.status_code, resp.content)
            )
            return False, "", 0
        uni_print("Resume uploading key <%s>" % key)
        return True, upload_id, int(resp["count"])

    @classmethod
    def init_multipart(cls, bucket, key):
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.initiate_multipart_upload(key)
        if resp.status_code == HTTP_OK:
            upload_id = resp["upload_id"]
        elif resp.status_code == HTTP_BAD_REQUEST:
            current_bucket.delete_object(key)
            resp = current_bucket.initiate_multipart_upload(key)
            upload_id = resp["upload_id"]
        else:
            upload_id = ""
        return upload_id

    @classmethod
    def upload_multipart(
            cls, filepath, upload_id, bucket, key, cur_part_number
    ):
        with open(filepath, "rb") as f:
            fc = FileChunk(f)
            if cur_part_number >= fc.parts:
                uni_print(
                    "Warning: The size of local file <%s> is smaller than previous "
                    "uploading size. Fail to resume previous uploading. Start a new"
                    "uploading." % filepath
                )
                cur_part_number = 0
            uni_print("File <%s> is uploading as Key <%s>" % (filepath, key))
            if cls.options.no_progress:
                pbar = None
            else:
                pbar = tqdm(
                    initial=cur_part_number * PART_SIZE,
                    total=fc.size,
                    unit="B",
                    unit_scale=True,
                    desc="Transferring",
                    bar_format=BAR_FORMAT,
                    leave=False,
                    ascii=USE_ASCII
                )
            for (part_number, data) in fc.iter(cur_part_number):
                is_upload_success = cls.upload_part(
                    upload_id,
                    part_number,
                    wrapper_stream(data, pbar, cls.get_tokens_obj()),
                    bucket,
                    key,
                )
                if not is_upload_success:
                    return (False, part_number)
            if pbar:
                pbar.close()
            return (True, fc.parts)

    @classmethod
    def upload_part(cls, upload_id, part_number, data, bucket, key):
        """upload one part of large file.
        """
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        resp = current_bucket.upload_multipart(
            key, upload_id=upload_id, part_number=part_number, body=data
        )
        if resp.status_code != HTTP_OK_CREATED:
            uni_print(resp.content)
            return False
        else:
            return True

    @classmethod
    def complete_multipart(cls, filepath, upload_id, cur_parts, bucket, key):
        parts = []
        for part_number in range(cur_parts):
            parts.append({"part_number": part_number})
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        cls.recorder.remove_record(filepath, bucket, key)
        resp = current_bucket.complete_multipart_upload(
            key, upload_id=upload_id, object_parts=parts
        )
        if resp.status_code != HTTP_OK_CREATED:
            uni_print(resp.content)
        else:
            statement = "Key <%s> created in bucket <%s>" % (key, bucket)
            uni_print(statement)
            if cls.command == "mv":
                os.remove(filepath)

    @classmethod
    def cleanup(cls, transfer_flow, bucket, prefix):
        pass

    @classmethod
    def send_request(cls):
        # if has option.limit_rate, create tokens object
        if hasattr(cls.options, "rate_limit"):
            if cls.options.rate_limit:
                cls.set_tokens_obj(cls.options.rate_limit)

        # Register SIGINT handler
        signal.signal(signal.SIGINT, cls._handle_sigint)
        transfer_flow = cls.get_transfer_flow()
        if transfer_flow == "LOCAL_TO_QS":
            if (cls.command == "sync") or (cls.options.recursive is True):
                cls.upload_files()
            else:
                cls.upload_file()
        elif transfer_flow == "QS_TO_LOCAL":
            if (cls.command == "sync") or (cls.options.recursive is True):
                cls.download_files()
            else:
                cls.download_file()
        if cls.recorder:
            cls.recorder.close()

    @classmethod
    def set_tokens_obj(cls, tokens_num, fill_rate=None):
        """
        :param tokens_num: the number of tokens
        :param fill_rate: the numbers of tokens that will fill in bucket in a second
        :return:
        """
        limit_rate = convert_to_bytes(tokens_num)
        # Set capacity=fill_rate *1.2 to avoid network jitter
        cls.tokens = TokenPail(capacity=limit_rate * 1.2, fill_rate=limit_rate)

    @classmethod
    def get_tokens_obj(cls):
        return cls.tokens

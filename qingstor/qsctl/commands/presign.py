# -*- coding: UTF-8 -*-
# =========================================================================
# Copyright (C) 2017 Yunify, Inc.
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

import sys

from .base import BaseCommand

from ..utils import get_current_time


class PresignCommand(BaseCommand):
    command = "presign"
    usage = "%(prog)s <qs-path> [-e <expire_seconds>]"

    @classmethod
    def add_extra_arguments(cls, parser):

        parser.add_argument(
            "qs_path",
            nargs="?",
            default="qs://",
            help="The qs-path to presign"
        )

        parser.add_argument(
            "-e",
            "--expire",
            dest="expire_seconds",
            type=int,
            default=3600,
            help="The number of seconds until the pre-signed URL expires."
        )
        return parser

    @classmethod
    def generate_presign_url(cls):
        bucket, prefix = cls.validate_qs_path(cls.options.qs_path)
        if prefix == "":
            cls.uni_print("Error: please specify object in qs-path")
            sys.exit(-1)
        resp = cls.current_bucket.head_object(prefix)

        # Handle common errors
        if resp.status_code == 404:
            cls.uni_print("Error: Please check if object <%s> exists" % prefix)
            sys.exit(-1)
        if resp.status_code == 403:
            cls.uni_print(
                "Error: Please check if you have enough"
                " permission to access object <%s>." % prefix
            )
            sys.exit(-1)
        if resp.status_code != 200:
            cls.uni_print(resp.content)
            sys.exit(-1)

        is_public = False
        # check whether the bucket is public
        current_acl = cls.current_bucket.get_acl()
        if current_acl.status_code == 200:
            for v in current_acl["acl"]:
                if v["grantee"]["name"] == "QS_ALL_USERS":
                    is_public = True

        if is_public:
            public_url = "{protocol}://{bucket_name}.{zone}.{host}/{object_key}".format(
                protocol=cls.current_bucket.config.protocol,
                bucket_name=bucket,
                zone=cls.bucket_map[bucket],
                host=cls.current_bucket.config.host,
                object_key=prefix
            )
            cls.uni_print(public_url)
            return public_url
        else:
            # if the bucket is non-public, generate the link with signature,
            # expire seconds and other formatted parameters
            prepared = cls.current_bucket.get_object_request(prefix).sign_query(
                get_current_time() + cls.options.expire_seconds
            )
            cls.uni_print(prepared.url)
            return prepared.url

    @classmethod
    def send_request(cls):
        cls.generate_presign_url()

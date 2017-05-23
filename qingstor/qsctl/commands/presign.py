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
from ..utils import get_current_time, uni_print


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
            uni_print("Error: please specify object in qs-path")
            sys.exit(-1)
        cls.validate_bucket(bucket)
        current_bucket = cls.client.Bucket(bucket, cls.bucket_map[bucket])
        prepared = current_bucket.get_object_request(prefix).sign_query(
            get_current_time() + cls.options.expire_seconds
        )
        uni_print(prepared.url)
        return prepared.url

    @classmethod
    def send_request(cls):
        cls.generate_presign_url()

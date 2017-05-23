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

import sys

from .base import BaseCommand
from ..utils import uni_print


class RmCommand(BaseCommand):

    command = "rm"
    usage = (
        "%(prog)s <qs_path> [-c <conf_file> -r <recursive> "
        "--exclude <pattern value> --include <pattern value>]"
    )

    @classmethod
    def add_extra_arguments(cls, parser):
        parser.add_argument(
            "qs_path", help="Key or keys under a specific prefix to be deleted"
        )

        parser.add_argument(
            "-r",
            "--recursive",
            action="store_true",
            dest="recursive",
            help="Recursively delete keys under a specific prefix"
        )

        parser.add_argument(
            "--exclude",
            type=str,
            help="Exclude all files or keys that match the specified pattern"
        )

        parser.add_argument(
            "--include",
            type=str,
            help="Do not exclude files or keys that match the specified pattern"
        )
        return parser

    @classmethod
    def send_request(cls):
        bucket, prefix = cls.validate_qs_path(cls.options.qs_path)
        if cls.options.recursive == True:
            cls.remove_multiple_keys(bucket, prefix)
        else:
            key = prefix
            if key == "":
                uni_print(
                    "Error: You must give a correct and complete qs_path, "
                    "such as 'qs://testbucket/testfile'."
                )
                sys.exit(-1)
            cls.remove_key(bucket, key)

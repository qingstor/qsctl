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

import sys

from .base import BaseCommand

from ..constants import HTTP_OK_CREATED

class MbCommand(BaseCommand):

    command = "mb"
    usage = "%(prog)s <bucket> [-c <conf_file> -z <zone>]"

    @classmethod
    def add_extra_arguments(cls, parser):
        parser.add_argument(
            "bucket",
            help="Name of the bucket to be created"
        )

        parser.add_argument(
            "-z",
            "--zone",
            dest="zone",
            help="In which zone to create the bucket"
        )
        return parser

    @classmethod
    def send_request(cls, options):
        bucket, prefix = cls.validate_qs_path(options.bucket)
        if prefix != "":
            print("Error: Invalid bucket name")
            sys.exit(-1)
        headers = {}
        if options.zone:
            headers["Location"] = options.zone
        resp = cls.conn.make_request("PUT", bucket, headers=headers)
        if resp.status == HTTP_OK_CREATED:
            print("Bucket <%s> created" % bucket)
        else:
            print(resp.read())
        resp.close()

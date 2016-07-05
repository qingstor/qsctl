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

from .transfer import TransferCommand

from ..utils import is_pattern_match, to_unix_path

class MvCommand(TransferCommand):

    command = "mv"
    usage = ("%(prog)s <source-path> <dest-path> [-c <conf_file> "
             "-r <recusively> --exclude <pattern value> --include <pattern value>]")

    @classmethod
    def cleanup(cls, transfer_method, options, bucket, prefix):
        if transfer_method == "PUT":
            cls.clean_files(options)

    @classmethod
    def clean_files(cls, options):
        for rt, dirs, files in os.walk(options.source_path):
            for d in dirs:
                local_path = os.path.join(rt, d)
                key_path = os.path.relpath(local_path, options.source_path) + "/"
                key_path = to_unix_path(key_path)
                if not os.listdir(local_path) and is_pattern_match(key_path, \
                    options.exclude, options.include):
                    os.rmdir(local_path)
                    print("Local directory '%s' deleted" % local_path)

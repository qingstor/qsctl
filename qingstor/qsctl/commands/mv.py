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

from .transfer import TransferCommand

from ..utils import is_pattern_match, to_unix_path, uni_print


class MvCommand(TransferCommand):

    command = "mv"
    usage = (
        "%(prog)s <source-path> <dest-path> [-c <conf_file> "
        "-r <recusively> --exclude <pattern value> --include <pattern value> --rate-limit <pattern value>]"
    )

    @classmethod
    def clean_empty_dirs(cls):
        local_dirs = []
        for rt, dirs, files in os.walk(cls.options.source_path):
            for d in dirs:
                local_dirs.append(os.path.join(rt, d))

        for local_dir in local_dirs[::-1]:
            key_path = os.path.relpath(local_dir, cls.options.source_path) + "/"
            key_path = to_unix_path(key_path)

            # Delete empty directory.
            if not os.listdir(local_dir) and is_pattern_match(
                    key_path, cls.options.exclude, cls.options.include
            ):
                os.rmdir(local_dir)
                uni_print("Local directory '%s' deleted" % local_dir)

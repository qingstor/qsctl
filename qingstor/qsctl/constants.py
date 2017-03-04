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

from .compat import is_python26, is_windows

# Buffer size
BUFFER_SIZE = 4 * 1024 * 1024

# Size of one part of the large file, used by Multipart-upload
PART_SIZE = 32 * 1024 * 1024

# HTTP response status
HTTP_OK = 200
HTTP_OK_CREATED = 201
HTTP_OK_NO_CONTENT = 204
HTTP_OK_PARTIAL_CONTENT = 206
HTTP_BAD_REQUEST = 400

# tqdm bar format
BAR_FORMAT = "{l_bar}{bar}| {n_fmt}/{total_fmt} [{remaining} ETA, {rate_fmt}]"

# Make tqdm work well on python2.6, windows and non utf-8 shell
USE_ASCII = is_python26 or is_windows or sys.stdin.encoding is not "UTF-8"

# Units used in output
UNITS = ('KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB')

# Temporary suffix for the file in downloading process.
TEMPORARY_FILE_SUFFIX = ".qsdownload"

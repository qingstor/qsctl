# -*- coding: utf-8 -*-
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
import platform

try:
    from yaml import CLoader as Loader
    from StringIO import cStringIO as StringIO
    import cPickle as pickle
except ImportError:
    from yaml import Loader
    from io import BytesIO as StringIO
    import pickle

_ver = sys.version_info
is_python2 = (_ver[0] == 2)
is_python3 = (_ver[0] == 3)
is_python26 = (_ver[0] == 2) and (_ver[1] == 6)

is_windows = platform.system().lower() == 'windows'

stdin_encoding = sys.stdin.encoding or "UTF-8"
stdout_encoding = sys.stdout.encoding or "UTF-8"

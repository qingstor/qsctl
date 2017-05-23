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

import re
import os
import json
import time
import calendar
from math import ceil
from yaml import load
from .constants import PART_SIZE, UNITS
from .compat import (
    is_python2, is_python3, is_windows, Loader, StringIO, stdout_encoding
)
from threading import RLock


class UploadIdRecorder(object):
    """
    This class stores upload_id for uploading a large object via multipart API.

    - In-memory record dict:
    key: <local_path>|<bucket>/<key> value: <upload_id>

    - On-disk record format:
    <local_path>|<bucket>/<key>#<upload_id>\n
    """

    def __init__(self, record_filename):
        self.separator = "|"
        self.records = {}
        self.dirty = False

        if os.path.exists(record_filename):
            self.file = open(record_filename, "r+")
            # Load records from file.
            records = self.file.readlines()
            for record in records:
                if is_python2:
                    record = record.decode("utf-8")
                kv = record.rsplit("#", 1)
                if len(kv) == 2:
                    # Remove the trailing \n
                    key = kv[1][:-1]
                    self.records[kv[0]] = key
        else:
            self.file = open(record_filename, "w+")

    def put_record(self, full_path, bucket, key, upload_id):
        key = self._get_record_key(full_path, bucket, key)
        self.records[key] = upload_id
        self.dirty = True

    def get_record(self, full_path, bucket, key):
        key = self._get_record_key(full_path, bucket, key)
        return self.records.get(key, "")

    def remove_record(self, full_path, bucket, key):
        key = self._get_record_key(full_path, bucket, key)
        self.records.pop(key, None)
        self.dirty = True

    def _get_record_key(self, local_path, bucket, key):
        full_path = os.path.join(os.getcwd(), local_path)
        return "%s%s%s/%s" % (full_path, self.separator, bucket, key)

    def close(self):
        if self.dirty:
            self._sync_record()
        self.file.close()

    def _sync_record(self):
        self.file.seek(0, 0)
        self.file.truncate()
        for key, value in self.records.items():
            record = "%s#%s\n" % (key, value)
            if is_python2:
                record = record.encode("utf-8")
            self.file.write(record)
        self.file.flush()


def yaml_load(stream):
    """
    Load from yaml stream and create a new python object

    @return object or `None` if failed
    """
    try:
        obj = load(stream, Loader=Loader)
    except Exception as e:
        print(e)
        obj = None
    return obj


def load_conf(conf_file):
    require_params = ["access_key_id", "secret_access_key"]
    compatible_params = ["qy_access_key_id", "qy_secret_access_key"]

    if conf_file == "":
        print("Config file should be specified")
        return None

    if conf_file.startswith('~'):
        conf_file = os.path.expanduser(conf_file)

    if not os.path.isfile(conf_file):
        print("Config file [%s] not exists" % conf_file)
        return None

    with open(conf_file, "r") as fd:
        conf = yaml_load(fd)
        if conf is None:
            print("Config file [%s] format error" % conf_file)
            return None
        for param in compatible_params:
            if param in conf:
                conf[param[3:]] = conf[param]
        for param in require_params:
            if param not in conf:
                print("[%s] should be specified in conf_file" % param)
                return None
    return conf


def confirm_by_user(notice):
    while True:
        inp = input(notice) if is_python3 else raw_input(
            notice.encode(stdout_encoding)
        )
        if inp == "y":
            return True
        if inp == "n":
            return False


def to_unix_path(path):
    if path is not None:
        path = path.replace("\\", "/")
    return path


def join_local_path(local_path, key_name):
    if is_windows:
        key_name = key_name.replace("/", "\\")
    local_path = os.path.join(local_path, key_name)
    return local_path


def uni_print(statement):
    """This function is used to properly write unicode to console.
    It ensures that the proper encoding is used in different os platforms.
    """
    try:
        if is_python2:
            statement = statement.encode(stdout_encoding)
        print(statement)
    except UnicodeError:
        print(
            "Warning: Your shell's encoding <%s> does not "
            "support printing this content" % stdout_encoding
        )


def json_loads(s):
    try:
        obj = json.loads(s)
    except:
        obj = json.loads(s.decode())
    return obj


def format_size(value):
    """Convert a size in number into: 'Byte', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB'.
    """
    base = 1024
    bytes_int = float(value)

    if bytes_int == 1:
        return '1 Byte'
    elif bytes_int < base:
        return '%d Bytes' % bytes_int

    for i, unit in enumerate(UNITS):
        unit_size = base**(i + 2)
        if round((bytes_int / unit_size) * base) < base:
            return '%.1f %s' % ((base * bytes_int / unit_size), unit)


def pattern_match(s, p):
    """pattern match used in 'include' and 'exclude' option
    """
    i, j, star_match_pos, last_star_pos = 0, 0, 0, -1
    while i < len(s):
        if j < len(p) and p[j] in (s[i], '?'):
            i, j = i + 1, j + 1
        elif j < len(p) and p[j] == '*':
            star_match_pos, last_star_pos = i, j
            j += 1
        elif last_star_pos > -1:
            i, star_match_pos = star_match_pos + 1, star_match_pos + 1
            j = last_star_pos + 1
        else:
            return False
    while j < len(p) and p[j] == '*':
        j += 1
    return j == len(p)


def is_pattern_match(s, exclude, include):
    """check if pattern match with 'include' and 'exclude' option
    """
    if is_windows:
        exclude = to_unix_path(exclude)
        include = to_unix_path(include)
    if exclude == None:
        return True
    elif include == None:
        return (not pattern_match(s, exclude))
    else:
        return (not pattern_match(s, exclude) or pattern_match(s, include))


class FileChunk:

    def __init__(self, fileobj):
        self.fileobj = fileobj

        # Handle files with do not support seek
        try:
            self.fileobj.seek(0)
        except IOError:
            # FIXME: this will read all data into memory
            self.fileobj = StringIO(self.fileobj.read())

        # Get file size
        self.fileobj.seek(0, os.SEEK_END)
        self.size = self.fileobj.tell()
        self.fileobj.seek(0, os.SEEK_SET)

        # Get parts
        self.parts = int(ceil(self.size * 1.0 / PART_SIZE))

    def iter(self, offset=0):
        for i in range(offset, self.parts):
            self.fileobj.seek(i * PART_SIZE)
            data = StringIO(self.fileobj.read(PART_SIZE))
            yield (i, data)


def wrapper_stream(stream, pbar=None, tokens=None):
    """
    Wrap stream.read() to upload progress bar
    """
    if not pbar:
        return stream

    _read = stream.read

    def _wrapper(size=None):
        buf = _read(size)
        pbar.update(size)
        return buf

    def _wrapper_with_rate_limit(size=None):
        # Use token bucket, wait while there are not enough tokens
        # If total tokens less than size, set size=total tokens
        if tokens.get_total_tokens() < size:
            size = tokens.get_total_tokens()
            size = int(size)
        while not tokens.consume(size):
            continue
        buf = _read(size)
        pbar.update(size)
        return buf

    _wrapper.__name__ = str("read")
    _wrapper_with_rate_limit.__name__ = str("read")
    if isinstance(tokens, TokenPail):
        # rate limit
        stream.read = _wrapper_with_rate_limit
    else:
        stream.read = _wrapper
    return stream


def validate_bucket_name(bucket_name):
    """
    Validate bucket name

    Bucket name must be compatible with DNS name (RFC 1123):

      - Less than 63 characters
      - Valid character set [a-z0-9-]
      - Can not begin and end with "-"

    Returns Trues if valid, False otherwise
    """
    if len(bucket_name) < 6 or len(bucket_name) > 63:
        return False

    if bucket_name.startswith("-") or bucket_name.endswith("-"):
        return False

    pattern = re.compile("^[0-9a-z]([0-9a-z-]{0,61})[0-9a-z]$")

    if not pattern.match(bucket_name):
        return False

    return True


def get_current_time():
    return calendar.timegm(time.gmtime())


# Implementation of token bucket algorithm that use to rate limit
class TokenPail(object):

    def __init__(self, capacity, fill_rate, is_lock=False):
        """
        :param capacity:  The total tokens in the bucket.
        :param fill_rate:  The rate in tokens/second that the bucket will be refilled
        """
        self._capacity = float(capacity)
        self._tokens = float(capacity)
        self._fill_rate = float(fill_rate)
        self._last_time = time.time()
        self._is_lock = is_lock
        self._lock = RLock()

    def _get_cur_tokens(self):
        if self._tokens < self._capacity:
            now = time.time()
            delta = self._fill_rate * (now - self._last_time)
            self._tokens = min(self._capacity, self._tokens + delta)
            self._last_time = now
        return self._tokens

    def get_cur_tokens(self):
        if self._is_lock:
            with self._lock:
                return self._get_cur_tokens()
        else:
            return self._get_cur_tokens()

    def _consume(self, tokens):
        if tokens <= self.get_cur_tokens():
            self._tokens -= tokens
            return True
        return False

    def consume(self, tokens):
        if self._is_lock:
            with self._lock:
                return self._consume(tokens)
        else:
            return self._consume(tokens)

    def get_total_tokens(self):
        return self._capacity


def convert_to_bytes(data):
    """
    Convert K/M/G to Bytes
    :param data: eg: 10K
    :return: bytes
    """
    multi = 1000
    unlimit = multi * multi * multi * multi
    result = unlimit
    if data:
        data = data.lower()
        try:
            if data.endswith("k"):
                data = data[:-1]
                result = int(data) * multi
            elif data.endswith("m"):
                data = data[:-1]
                result = int(data) * multi * multi
            elif data.endswith("g"):
                data = data[:-1]
                result = int(data) * multi * multi * multi
            else:
                result = int(data)
        except ValueError:
            uni_print(
                "Warning: rate limit include invaild character," \
                "use 1G/s rate limit  as default"
            )
            result = unlimit
    if result <= 0:
        uni_print(
            "Warning: rate limit cannot be negative," \
            "use 1G/s rate limit  as default"
        )
        result = unlimit
    return result

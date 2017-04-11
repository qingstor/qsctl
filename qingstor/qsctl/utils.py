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
import sys
import json
import time
import calendar
from yaml import load
from .constants import PART_SIZE, UNITS
from .compat import (
    is_python2, is_python3, is_windows, Loader, StringIO, stdout_encoding
)


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
    '''
    Load from yaml stream and create a new python object

    @return object or `None` if failed
    '''
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
        inp = input(notice) if is_python3 else raw_input(notice.encode(stdout_encoding))
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
    '''pattern match used in 'include' and 'exclude' option
    '''
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
    '''check if pattern match with 'include' and 'exclude' option
    '''
    if is_windows:
        exclude = to_unix_path(exclude)
        include = to_unix_path(include)
    if exclude == None:
        return True
    elif include == None:
        return (not pattern_match(s, exclude))
    else:
        return (not pattern_match(s, exclude) or pattern_match(s, include))


def get_part_numbers(filename):
    '''return a list of part numbers, will be used in multipart upload.
    '''
    part_numbers = []
    filesize = os.path.getsize(filename)
    num = filesize // PART_SIZE
    if filesize % PART_SIZE != 0:
        num = num + 1
    for i in range(0, num):
        part_numbers.append(i)
    return part_numbers


class FileChunk(object):

    def __init__(self, filepath, part_number):
        self._filepath = filepath
        self._start_byte = PART_SIZE * part_number
        self._fileobj = open(self._filepath, 'rb')
        self._size = self._calculate_chunk_size(self._fileobj, self._start_byte)
        self._fileobj.seek(self._start_byte)
        self._amount_read = 0

    def _calculate_chunk_size(self, fileobj, start_byte):
        actual_file_size = os.fstat(fileobj.fileno()).st_size
        max_chunk_size = actual_file_size - start_byte
        return min(max_chunk_size, PART_SIZE)

    def read(self, amount=None):
        if amount is None:
            remaining = self._size - self._amount_read
            data = self._fileobj.read(remaining)
            self._amount_read += remaining
            return data
        else:
            actual_amount = min(self._size - self._amount_read, amount)
            data = self._fileobj.read(actual_amount)
            self._amount_read += actual_amount
            return data

    def seek(self, where):
        self._fileobj.seek(self._start_byte + where)
        self._amount_read = where

    def close(self):
        self._fileobj.close()

    def tell(self):
        return self._amount_read

    def __len__(self):
        '''__len__ is defined because requests will try to determine the length
        of the stream to set a content length.
        '''
        return self._size

    def __enter__(self):
        return self

    def __exit__(self, *args, **kwargs):
        self._fileobj.close()

    def __iter__(self):
        '''Basically httplib will try to iterate over the contents, even
        if it is a file like object.
        '''
        return iter([])


class StdinFileChunk(StringIO):

    def __init__(self, max_size):
        StringIO.__init__(self)
        self.write(sys.stdin.read(max_size))
        self.seek(0, os.SEEK_SET)

    def __len__(self):
        pos = self.tell()
        self.seek(0, os.SEEK_END)
        l = self.tell()
        self.seek(pos, os.SEEK_SET)
        return l


def wrapper_stream(stream, pbar=None):
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

    _wrapper.__name__ = str("read")
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

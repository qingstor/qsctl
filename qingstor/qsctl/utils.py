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
import sys
import json
import platform
from yaml import load, dump

try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

from .constants import PART_SIZE

UNITS = ('KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB')

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
    require_params = [
        "qy_access_key_id",
        "qy_secret_access_key",
        "zone",
    ]

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
        for param in require_params:
            if param not in conf:
                print("[%s] should be specified in conf_file" % param)
                return None
    return conf

def confirm_by_user(notice):
    while True:
        inp = input(notice) if sys.version > "3" else raw_input(notice)
        if inp == "y":
            return True
        if inp == "n":
            return False

def to_unix_path(path):
    if path is not None:
        if is_windows() and sys.version < "3":
            path = encode_to_utf8(path)
        path = path.replace("\\","/")
    return path

def join_local_path(local_path, key_name):
    if is_windows():
        if sys.version < "3":
            key_name = encode_to_gbk(key_name)
        key_name = key_name.replace("/","\\")
    local_path = os.path.join(local_path, key_name)
    return local_path

def encode_to_utf8(s):
    return s.decode('gbk').encode('utf8')

def encode_to_gbk(s):
    return s.decode('utf8').encode('gbk')

def uni_print(statement):
    """This function is used to properly write unicode to console.
    It ensures that the proper encoding is used in different os platforms.
    """
    if is_windows() and sys.version < "3":
        statement = statement.decode('utf8')
    print(statement)

def is_windows():
    return platform.system().lower() == 'windows'

def json_loads(s):
    try:
        obj = json.loads(s)
    except:
        obj = json.loads(s.decode())
    return obj

def format_size(value):
    """Convert a size in number into: 'Byte', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB'.
    """
    one_decimal_point = '%.1f'
    base = 1024
    bytes_int = float(value)

    if bytes_int == 1:
        return '1 Byte'
    elif bytes_int < base:
        return '%d Bytes' % bytes_int

    for i, unit in enumerate(UNITS):
        unit_size = base ** (i+2)
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
    while j < len(p) and p[j] == '*': j += 1
    return j == len(p)

def is_pattern_match(s, exclude, include):
    '''check if pattern match with 'include' and 'exclude' option
    '''
    if is_windows():
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

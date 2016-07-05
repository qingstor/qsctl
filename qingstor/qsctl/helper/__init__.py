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
from subprocess import Popen, PIPE

from docutils.core import publish_string
from docutils.writers import manpage

from qingstor.qsctl.utils import is_windows
from qingstor.qsctl.helper.textwriter import TextWriter
from qingstor.qsctl.helper.qsdocutils import (
    ignore_ctrl_c,
    gen_rst_doc
)

def get_renderer(command):
    """
    Return appropriate HelpRenderer for different platform.
    """
    if is_windows():
        return WindowsHelpRenderer(command)
    else:
        return PosixHelpRenderer(command)

class HelpRenderer(object):

    def __init__(self, command):
        self.command = command

    def _build_usage(self, command):
        raise NotImplementedError

    def _print_usage(self, usage):
        raise NotImplementedError

    def render(self):
        usage = self._build_usage(self.command)
        self._print_usage(usage)

class PosixHelpRenderer(HelpRenderer):

    def _build_usage(self, command):
        rst_doc = gen_rst_doc(command)
        man_doc = publish_string(rst_doc, writer=manpage.Writer())
        cmdline = ['groff', '-m', 'man', '-T', 'ascii']
        p = Popen(cmdline, stdin=PIPE, stdout=PIPE, stderr=PIPE)
        usage = p.communicate(input=man_doc)[0]
        return usage

    def _print_usage(self, usage):
        cmdline = ['less', '-r']
        with ignore_ctrl_c():
            # The default behavior of less ignore ctrl-c (you can't ctrl-c
            # out of a manpage).
            p = Popen(cmdline, stdin=PIPE)
            p.communicate(input=usage)

class WindowsHelpRenderer(HelpRenderer):

    def _build_usage(self, command):
        rst_doc = gen_rst_doc(command)
        text_doc = publish_string(rst_doc, writer=TextWriter())
        return text_doc

    def _print_usage(self, usage):
        cmdline = ['more']
        p = Popen(cmdline, stdin=PIPE, shell=True)
        p.communicate(input=usage)

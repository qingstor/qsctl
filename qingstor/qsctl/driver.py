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
import argparse
from difflib import get_close_matches

from qingstor.qsctl import __version__
from qingstor.qsctl.utils import is_python2
from qingstor.qsctl.helper import get_renderer
from qingstor.qsctl.compat import stdin_encoding
from qingstor.qsctl.commands.ls import LsCommand
from qingstor.qsctl.commands.cp import CpCommand
from qingstor.qsctl.commands.rm import RmCommand
from qingstor.qsctl.commands.rb import RbCommand
from qingstor.qsctl.commands.mb import MbCommand
from qingstor.qsctl.commands.mv import MvCommand
from qingstor.qsctl.commands.sync import SyncCommand
from qingstor.qsctl.commands.presign import PresignCommand

COMMANDS = ('ls', 'cp', 'mb', 'mv', 'rb', 'rm', 'sync', 'presign')

INDENT = ' ' * 2
NEWLINE = '\n' + INDENT


def exit_due_to_invalid_command(suggest_commands=None):
    usage = NEWLINE + '%(prog)s <command> [parameters]\n\n' \
        + 'To see help text, you can run:\n\n' \
        + INDENT + '%(prog)s help\n' \
        + INDENT + '%(prog)s <command> help\n\n' \
        + 'Valid commands are:\n\n' \
        + INDENT + NEWLINE.join(COMMANDS)

    if suggest_commands:
        usage += '\n\nInvalid command, you might want:\n\n' \
            + ', '.join(suggest_commands)

    parser = argparse.ArgumentParser(
        prog='qsctl',
        usage=usage,
    )
    parser.print_help()
    sys.exit(-1)


def check_argument(args):
    if is_python2:
        for i in range(len(args)):
            args[i] = args[i].decode(stdin_encoding)

    if len(args) < 2:
        exit_due_to_invalid_command()

    if args[1].lower() in ('-v', '-version', 'v', 'version'):
        print('qsctl %s' % __version__)
        sys.exit(0)

    if args[-1].lower() == 'help' and len(args) <= 3:
        args[0] = "qsctl"
        command = "-".join(args[:-1])
        renderer = get_renderer(command)
        renderer.render()
        sys.exit(0)

    command = args[1]

    if command not in COMMANDS:
        suggest_commands = get_close_matches(command, COMMANDS)
        exit_due_to_invalid_command(suggest_commands)


def get_command(command):
    if command in COMMANDS:
        return globals()[command.capitalize() + "Command"]
    else:
        exit_due_to_invalid_command()


def main():
    args = sys.argv
    check_argument(args)
    command = get_command(args[1])
    command.main(args[2:])

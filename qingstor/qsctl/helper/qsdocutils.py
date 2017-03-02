# -*- coding: UTF-8 -*-

import os
import signal
import contextlib

# Used to build usages
AUTHOR = "qs-devel@yunify.com"
COPYRIGHT = "Copyright (C) 2016-2017 Yunify, Inc"

COMMANDS = (
    'qsctl', 'qsctl-ls', 'qsctl-cp', 'qsctl-mb', 'qsctl-mv', 'qsctl-rb',
    'qsctl-rm', 'qsctl-sync', 'qsctl-presign'
)


def to_rst_style_title(title):
    style = "=" * len(title)
    return "%s\n%s\n%s" % (style, title, style)


def gen_see_also(command):
    content = ""
    for c in COMMANDS:
        if c != command:
            content += "* ``%s help``\n\n" % c.replace("-", " ")
    return content


class RstDocument(object):

    def __init__(self, rst_source=""):
        self.rst_source = rst_source

    def from_file(self, filepath):
        with open(filepath) as f:
            self.rst_source += f.read()

    def add_reporting_bug(self):
        title = to_rst_style_title("Reporting Bug")
        content = "Report bugs to email <%s>." % AUTHOR
        block = "\n%s\n\n%s\n" % (title, content)
        self.rst_source += block

    def add_see_also(self, command):
        title = to_rst_style_title("See Also")
        content = gen_see_also(command)
        block = "\n%s\n\n%s\n" % (title, content)
        self.rst_source += block

    def add_copyright(self):
        title = to_rst_style_title("Copyright")
        content = COPYRIGHT
        block = "\n%s\n\n%s\n" % (title, content)
        self.rst_source += block

    def getvalue(self):
        return self.rst_source


def gen_rst_doc(command):
    rst_doc = RstDocument()

    # Descriptions and examples are in the 'source' directory.
    # We need read them out first.
    current_path = os.path.split(os.path.realpath(__file__))[0]
    source_path = "source/%s.rst" % command
    filepath = os.path.join(current_path, source_path)

    rst_doc.from_file(filepath)
    rst_doc.add_reporting_bug()
    rst_doc.add_see_also(command)
    rst_doc.add_copyright()
    return rst_doc.getvalue()


def gen_sphinx_doc(command):
    # Generating ReST documents for sphinx
    sphinx_doc = RstDocument()

    # Descriptions and examples are in the 'source' directory.
    # We need read them out first.
    current_path = os.path.split(os.path.realpath(__file__))[0]
    source_path = "source/%s.rst" % command
    filepath = os.path.join(current_path, source_path)

    sphinx_doc.from_file(filepath)
    return sphinx_doc.getvalue()


@contextlib.contextmanager
def ignore_ctrl_c():
    # Ctrl-c shoule be ignored when using 'less -r' to print usage.
    original = signal.signal(signal.SIGINT, signal.SIG_IGN)
    try:
        yield
    finally:
        signal.signal(signal.SIGINT, original)

#!/usr/bin/env python

import os
from qingstor.qsctl.helper.qsdocutils import (
    gen_sphinx_doc,
    COMMANDS,
)

REF_PATH = 'reference'

if not os.path.isdir(REF_PATH):
    os.mkdir(REF_PATH)

print('Generating ReST documents for all commands...')
for command in COMMANDS:
    contents = gen_sphinx_doc(command)
    filepath = "%s/%s.rst" % (REF_PATH, command)
    with open(filepath, 'wb') as f:
        f.write(contents)
print('Done!')

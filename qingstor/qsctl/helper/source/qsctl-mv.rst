.. _qsctl-mv:


********
qsctl-mv
********

========
Synopsis
========

::

      qsctl mv
    <local-path> <qs-path> or <qs-path> <local-path>
    [--force]
    [-r, --recursive]
    [--exclude]
    [--include]

===========
Description
===========

Move local file(s) to qingstor or qingstor key(s) to local. The first path
argument is the source and second the destination.

``mv`` command will delete source files(keys) after transfering work done.
If you do not want to delete source files(keys), please use ``cp`` command.

When a key or file already exists at destination, program will ask to
overwrite or skip. You can use ``--force`` option to forcely overwrite.

=======
Options
=======

``-f, --force``

Forcely overwrite existing key(file) without asking.

``-r, --recursive``

Recursively transfer keys(files).

``--exclude``

Exclude all keys or files that match the specified pattern.

``--include``

Do not exclude keys or files that match the specified pattern.

========
Examples
========

The following ``mv`` command moves a single file to bucket ``mybucket``::

    $ qsctl mv test.txt qs://mybucket

Output::

    Key <test.txt> created
    File 'test.txt' deleted

The following ``mv`` command moves all keys in bucket ``mybucket`` to local
directory::

    $ qsctl mv qs://mybucket test/ -r

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.jpg' written
    Key <test/test1.txt> deleted
    Key <test/test2.txt> deleted
    Key <test/test3.jpg> deleted

The following ``mv`` command with the ``--exclude`` option will exclude the keys
match the pattern value ``"*.txt"``::

    $ qsctl mv qs://mybucket test -r --exclude "*.txt"

Output::

    File 'test/test3.jpg' written
    Key <test/test3.jpg> deleted

The following ``mv`` command with the ``--exclude`` and ``--include`` option
will include the keys match the pattern value ``"*.txt"``::

    $ qsctl mv qs://mybucket test -r --exclude "*" --include "*.txt"

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    Key <test/test1.txt> deleted
    Key <test/test2.txt> deleted

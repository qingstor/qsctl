.. _qsctl-cp:


********
qsctl-cp
********


========
Synopsis
========

::

      qsctl cp
    <local-path> <qs-path> or <qs-path> <local-path>
    [--force]
    [-r, --recursive]
    [--exclude]
    [--include]
    [--rate-limit]

===========
Description
===========

Copy local file(s) to qingstor or qingstor key(s) to local. The first
path argument is the source and second the destination.

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

``--rate-limit``

Limit rate when cp file from local to qingstor, or qingstor to local,
unit: K/M/G, eg: 100K.

========
Examples
========

The following ``cp`` command copies a single file to bucket ``mybucket``::

    $ qsctl cp test.txt qs://mybucket

Output::

    Key <test.txt> created
    File 'test.txt' deleted

The following ``cp`` command copies all keys in bucket ``mybucket`` to local
directory::

    $ qsctl cp qs://mybucket test/ -r

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.jpg' written

The following ``cp`` command with the ``--exclude`` option will exclude the
keys match the pattern value ``"*.txt"`` ::

    $ qsctl cp qs://mybucket test -r --exclude "*.txt"

Output::

    File 'test/test3.jpg' written

The following ``cp`` command with the ``--exclude`` and ``--include`` option
will include the keys match the pattern value ``"*.txt"``::

    $ qsctl cp qs://mybucket test -r --exclude "*" --include "*.txt"

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

The following ``cp`` command copies a single file to bucket ``mybucket``,
and limit the transmission speed of 100K per second::

    $ qsctl cp test.txt qs://mybucket --rate-limit 100K

Output::

    Key <test.txt> created
    File 'test.txt' deleted

The following ``cp`` command copies test1.txt in bucket ``mybucket`` to local
directory, and limit the transmission speed of 100K per second::

    $ qsctl cp qs://mybucket/test1.txt test/ --rate-limit 100K

Output::

    File 'test/test1.txt' written

.. _qsctl-sync:


**********
qsctl-sync
**********


========
Synopsis
========

::

      qsctl sync
    <local-path> <qs-path> or <qs-path> <local-path>
    [--delete]
    [--exclude]
    [--include]
    [--rate-limit]

===========
Description
===========

Sync between local directory and QS-Directory. The first path argument is the
source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will
compare the modified time of source file(key) and destination key(file).
The destination key(file) will be overwritten only if the source one
newer than destination one.

=======
Options
=======

``--delete``

Any files or keys existing under the destination directory but not existing in
the source directory will be deleted if ``--delete`` option is specified.

``--exclude``

Exclude all keys or files that match the specified pattern.

``--include``

Do not exclude keys or files that match the specified pattern.

``--rate-limit``

Limit rate when sync file from local to qingstor, or qingstor to local,
unit: K/M/G, eg: 100K.

========
Examples
========

The following ``sync`` command will sync local directory to QS-Directory::

    $ qsctl sync . qs://mybucket

Output::

    Key <test1.txt> created
    Key <test2.txt> created
    Key <test3.txt> created

The following ``sync`` command sync QS-Directory to local directory::

    $ qsctl sync qs://mybucket/test test/

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

The following ``sync`` command with the ``--delete`` option will delete those
files existing in local directory but not existing in bucket ``mybucket``.
Suppose the local directory has file ``test3.txt``. The bucket ``mybucket``
only have ``test1.txt`` and ``test2.txt``::

    $ qsctl sync qs://mybucket/test/ test/ --delete

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.txt' deleted

The following ``sync`` command will sync local directory to QS-Directory,
and limit the transmission speed of 100K per second::

    $ qsctl sync . qs://mybucket --rate-limit 100K

Output::

    Key <test1.txt> created
    Key <test2.txt> created
    Key <test3.txt> created

The following ``sync`` command sync QS-Directory to local directory,
and limit the transmission speed of 100K per second::

    $ qsctl sync qs://mybucket/test test/ --rate-limit 100K

Output::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

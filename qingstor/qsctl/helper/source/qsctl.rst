.. _qsctl:


*****
qsctl
*****


========
Synopsis
========

::

    qsctl <command> [<args>] [<options>]

==============
Common Options
==============

``-c, --config``

Use a specific config file.

``-h, --help``

Prints the synopsis and available commands.

``-v``

Prints the qsctl version.

===========
Description
===========

Qsctl provides more powerful unix-like commands. You can manage QingStor
resources just like files on local machine. Unix-like commands contains:
cp, ls, mb, mv, rm, rb, and sync. All of them support batch processing.

This section also explains some important concepts and notations.

Path Argument Type
++++++++++++++++++

There are two types of path arguments: ``local-path`` and ``qs-path``.

``local-path``: represents the path of a local file or directory.

``qs-path``: represents the location of a qs key, prefix, or bucket. This
must be written in the form ``qs://mybucket/mykey`` where ``mybucket`` is
the specified qs bucket, ``mykey`` is the specified qs key. The path
argument must begin with ``qs://`` in order to denote that the path argument
refers to qs-path. Prefix is used for batch processing. For example, if you
run ``qsctl rm -r qs://mybucket/myprefix``, all keys under the specified
prefix ``myprefix`` in ``mybucket`` will be deleted. If prefix is empty(such
as ``qs://mybucket/``), the batch processing refers to all keys in the given
bucket.

QS-Directory
++++++++++++

``cp, mv`` comands with ``-r`` option and ``sync`` command perform only on
directories. In Qsctl, prefix ending with '/' will be seen as a ``qs-directory``.
In these directories operations, if a prefix not ending with '/' is given,
program will automatically add '/' at the end of the prefix. It means,
``qsctl sync /foo/bar/ qs://mybucket/myprefix`` will be executed as
``qsctl sync /foo/bar/ qs://mybucket/myprefix/``. Note that qs-path with
empty prefix is also a qs-directory(for exmaple ``qs://mybucket`` and
``qs://mybucket/``). It can be seen as the root directory of given bucket.

Order of Path Arguments
+++++++++++++++++++++++

All commands *excepting ls command* take one or two positional path arguments.
The first path argument represents the source. If there is a second path
argument it represents the destination. Commands with only one path argument
do not have a destination because the operation is being performed only on
the source.

Use of Exclude and Include Filters
++++++++++++++++++++++++++++++++++

Currently, there is no support for the use of UNIX style wildcards in path
arguments(such as on UNIX: ``cp -r dir1/* dir2/``). However, most commands
can achieve the desired result by adding ``--exclude "<value>"`` and
``--include "<value>"`` options.  These parameters perform pattern matching
to either exclude or include particular files or keys. The following pattern
symbols are supported.

* ``*``: Matches any characters
* ``?``: Matches any single character

For example, suppose you had the following directory structure::

    /tmp/foo/
      bar/
      |---mykey1
      |---mykey2
      test1.txt
      test2.txt
      test3.jpg

Given the directory structure above and the command
``qsctl cp /tmp/foo qs://bucket -r --exclude "bar/*"``, the files
``bar/mykey1`` and ``bar/mykey2`` will be excluded from the files to upload
because the exclude filter ``bar/*`` will have the source prepended to the
filter.  This means that::

    /tmp/foo/bar/* \-> /tmp/foo/bar/mykey1  (matches, should exclude)
    /tmp/foo/bar/* \-> /tmp/foo/bar/mykey2  (matches, should exclude)
    /tmp/foo/bar/* \-> /tmp/foo/test1.txt  (does not match, should include)
    /tmp/foo/bar/* \-> /tmp/foo/test2.txt  (does not match, should include)
    /tmp/foo/bar/* \-> /tmp/foo/test3.jpg  (does not match, should include)

Note that, by default, *all files are included*.  This means that providing
**only** an ``--include`` filter will not change what files are transferred.
``--include`` will only re-include files that have been excluded from an
``--exclude`` filter.  If you only want to upload files with a particular
extension, you need to first exclude all files, then re-include the files
with the particular extension. This command will upload **only** files ending
with ``.jpg``::

    $ qsctl cp /tmp/foo/ qs://mybucket/ -r --exclude "*" --include "*.jpg"

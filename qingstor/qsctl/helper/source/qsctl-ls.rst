.. _qsctl-ls:


********
qsctl-ls
********


========
Synopsis
========

::

      qsctl ls
    <qs-path> or none
    [-r, --recursive]
    [--page]
    [--zone]

===========
Description
===========

List qingstor keys under a prefix or all qingstor buckets.

=======
Options
=======

``-r, --recursive``

Recursively lists qingstor keys under a prefix.

``--page``

The number of results to return in each response.

``--zone``

List buckets located in this zone.

========
Examples
========

The following ``ls`` command lists all qingstor buckets::

    $ qsctl ls

Output::

    mybucket1
    mybucekt2

The following ``ls`` command lists all keys and qs-directory in the bucket
``mybucket``::

    $ qsctl ls qs://mybucket

Output::

    Directory                          myprefix/
    2016-04-03 11:16:04     4 Bytes    test1.txt
    2016-04-03 11:16:04     4 Bytes    test2.txt
    2016-04-03 11:16:04     4 Bytes    test3.txt

The following ``ls`` command uses the ``-r`` option to recursively list
all keys under prefix ``myprefix`` in bucket ``mybucket``::

    $ qsctl ls qs://mybucket/myprefix/ -r

Output::

    2016-04-03 11:16:04     4 Bytes    myprefix/test.txt
    2016-04-03 17:51:18     1.4 KiB    myprefix/test/test.txt
    2016-04-03 11:16:04     1.4 KiB    myprefix/test/test/test.txt


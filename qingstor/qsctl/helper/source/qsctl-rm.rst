.. _qsctl-rm:


********
qsctl-rm
********


========
Synopsis
========

::

      qsctl rm
    <qs-path>
    [-r, --recursive]
    [--exclude]
    [--include]

===========
Description
===========

Delete a qingstor key or keys under a prefix.

=======
Options
=======

``-r, --recursive``

Recursively delete qingstor keys under a prefix.

``--exclude``

Exclude all keys that match the specified pattern.

``--include``

Do not exclude keys that match the specified pattern.

========
Examples
========

The following ``rm`` command removes a qingstor key::

    $ qsctl rm qs://mybucket/test.txt

Output::

    Key <test.txt> deleted

The following ``rm`` command uses the ``-r`` option to recursively remove
all of the keys under the prefix ``myprefix`` in the bucket ``mybucket``::

    $ qsctl rm qs://mybucket/myprefix -r

Output::

     Key <myprefix.txt> deleted
     Key <myprefix/txt/test.txt> deleted
     Key <myprefix/jpg/test.jpg> deleted

The following ``rm`` command uses the ``--exclude`` option to exclude keys
match a specified pattern. Key ``<myprefix/test.jpg>`` is skipped::

    $ qsctl rm qs://mybucket/myprefix -r --exclude "*.jpg"

Output::

    Key <myprefix.txt> deleted
    Key <myprefix/test.txt> deleted

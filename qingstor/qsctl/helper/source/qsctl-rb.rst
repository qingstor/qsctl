.. _qsctl-rb:


********
qsctl-rb
********


========
Synopsis
========

::

      qsctl rb
    <bucket>
    [--force]

===========
Description
===========

Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.

=======
Options
=======

``--force``

Forcely delete nonempty qingstor bucket.

========
Examples
========

The following ``rb`` command removes an empty bucket::

    $ qsctl rb qs://mybucket

Output::

    Bucket <mybucket> deleted

The following ``rb`` command uses the --force option to first remove all of
the keys in the bucket and then remove the bucket itself::

    $ qsctl rb qs://mybucket --force

Output::

    Key <test1.txt> deleted
    Key <test2.txt> deleted
    Bucket <mybucket> deleted

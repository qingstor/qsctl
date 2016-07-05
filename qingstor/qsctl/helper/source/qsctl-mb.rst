.. _qsctl-mb:


********
qsctl-mb
********


========
Synopsis
========

::

      qsctl mb
    <bucket>
    [--zone]

===========
Description
===========

Create a QingStor bucket.

=======
Options
=======

``--zone``

Creates a bucket in a specified zone.

========
Examples
========

The following ``mb`` command creates a bucket. The bucket is created in the
zone specified in the user's configuration file::

    $ qsctl mb qs://mybucket

Output::

    Bucket <mybucket> created

The following mb command creates a bucket in a zone specified by the
``--zone`` option. In this example, the user makes the bucket ``mybucket``
in the zone ``pek3a``::

    $ qsctl mb qs://mybucket --zone pek3a

Output::

    Bucket <mybucket> created

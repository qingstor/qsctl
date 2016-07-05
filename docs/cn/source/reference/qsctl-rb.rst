.. _qsctl-rb:


********
qsctl-rb
********


====
总览
====

::

      qsctl rb
    <bucket>
    [--force]

====
描述
====

删除一个空的存储空间，或强制删除一个非空的存储空间。

=======
Options
=======

``--force``

强制删除一个非空的存储空间。

========
Examples
========

下面的 ``rb`` 命令删除一个空的存储空间::

    $ qsctl rb qs://mybucket

输出::

    Bucket <mybucket> deleted

下面的 ``rb`` 命令强制删除一个非空的存储空间::

    $ qsctl rb qs://mybucket --force

输出::

    Key <test1.txt> deleted
    Key <test2.txt> deleted
    Bucket <mybucket> deleted

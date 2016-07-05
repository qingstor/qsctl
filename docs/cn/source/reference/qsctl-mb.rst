.. _qsctl-mb:


********
qsctl-mb
********


====
总览
====

::

      qsctl mb
    <bucket>
    [--zone]

====
描述
====

创建一个新的存储空间。

====
选项
====

``--zone``

在给定的区域创建一个新的存储空间。

====
示例
====

下面的 ``mb`` 命令创建一个空的存储空间，区域由默认配置文件给定::

    $ qsctl mb qs://mybucket

输出::

    Bucket <mybucket> created

下面的 ``mb`` 命令创建一个空的存储空间，区域由 ``--zone`` 选项给定::

    $ qsctl mb qs://mybucket --zone pek3a

输出::

    Bucket <mybucket> created

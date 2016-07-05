.. _qsctl-ls:


********
qsctl-ls
********


====
总览
====

::

      qsctl ls
    <qs-path> or none
    [-r, --recursive]
    [--page]
    [--zone]

====
描述
====

列出给定前缀下的所有对象或列出所有存储空间。

====
选项
====

``-r, --recursive``

递归地列出给定前缀下的所有对象。

``--page``

每次http请求返回的对象数目。

``--zone``

列出某个区域的存储空间。

====
示例
====

下面的 ``ls`` 命令列出所有存储空间::

    $ qsctl ls

输出::

    mybucket1
    mybucekt2

下面的 ``ls`` 命令列出 ``mybucket`` 下的所有对象和QS-Direcory::

    $ qsctl ls qs://mybucket

输出::

    Directory                          myprefix/
    2016-04-03 11:16:04     4 Bytes    test1.txt
    2016-04-03 11:16:04     4 Bytes    test2.txt
    2016-04-03 11:16:04     4 Bytes    test3.txt

下面的 ``ls`` 命令列出所有前缀为 ``myprefix`` 的对象::

    $ qsctl ls qs://mybucket/myprefix/ -r

输出::

    2016-04-03 11:16:04     4 Bytes    myprefix/test.txt
    2016-04-03 17:51:18     1.4 KiB    myprefix/test/test.txt
    2016-04-03 11:16:04     1.4 KiB    myprefix/test/test/test.txt

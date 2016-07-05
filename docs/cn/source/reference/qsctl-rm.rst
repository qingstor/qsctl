.. _qsctl-rm:


********
qsctl-rm
********


====
总览
====

::

      qsctl rm
    <qs-path>
    [-r, --recursive]
    [--exclude]
    [--include]


====
描述
====

删除一个QingStor对象或给定前缀下的所有对象

====
选项
====

``-r, --recursive``

删除给定前缀下的所有对象。

``--exclude``

排除匹配特定类型的文件（对象）。

``--include``

包含匹配特定类型的文件（对象）。

====
示例
====

下面的 ``rm`` 命令删除一个QingStor对象::

    $ qsctl rm qs://mybucket/test.txt

输出::

    Key <test.txt> deleted

下面的 ``rm`` 命令删除所有前缀为 ``myprefix`` 的QingStor对象::

    $ qsctl rm qs://mybucket/myprefix -r

输出::

     Key <myprefix.txt> deleted
     Key <myprefix/txt/test.txt> deleted
     Key <myprefix/jpg/test.jpg> deleted

下面的 ``rm`` 命令将排除匹配类型 ``*.jpg`` 的QingStor对象::

    $ qsctl rm qs://mybucket/myprefix -r --exclude "*.jpg"

输出::

    Key <myprefix.txt> deleted
    Key <myprefix/test.txt> deleted

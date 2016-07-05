.. _qsctl-cp:


********
qsctl-cp
********


====
总览
====

::

      qsctl cp
    <local-path> <qs-path> or <qs-path> <local-path>
    [--force]
    [-r, --recursive]
    [--exclude]
    [--include]

====
描述
====

复制文件到QingStor或复制QingStor对象到本地。第一个Path变量是源地址，第二个
是目标地址。

如果要传的文件（对象）目标地址已存在，程序将会询问覆盖或跳过。你可以使用
``--force`` 选项强制覆盖。

====
选项
====

``-f, --force``

强制覆盖已存在的目标文件（对象）。

``-r, --recursive``

递归地传输目录下所有文件。

``--exclude``

排除匹配特定类型的文件（对象）。

``--include``

包含匹配特定类型的文件（对象）。

====
示例
====

下面的 ``cp`` 命令将一个本地文件复制到 ``mybucket``::

    $ qsctl cp test.txt qs://mybucket

输出::

    Key <test.txt> created
    File 'test.txt' deleted

下面的 ``cp`` 命令将 ``mybucket`` 下的所有对象复制到本地::

    $ qsctl cp qs://mybucket test/ -r

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.jpg' written

下面的 ``cp`` 命令将排除所有匹配 ``"*.txt"`` 的对象::

    $ qsctl cp qs://mybucket test -r --exclude "*.txt"

输出::

    File 'test/test3.jpg' written

下面的 ``cp`` 命令只包含匹配 ``"*.txt"`` 的对象::

    $ qsctl cp qs://mybucket test -r --exclude "*" --include "*.txt"

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

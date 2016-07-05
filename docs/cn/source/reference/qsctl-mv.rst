.. _qsctl-mv:


********
qsctl-mv
********


====
总览
====

::

      qsctl mv
    <local-path> <qs-path> or <qs-path> <local-path>
    [--force]
    [-r, --recursive]
    [--exclude]
    [--include]

====
描述
====

移动文件到QingStor或移动QingStor对象到本地。第一个Path变量是源地址，第二个
是目标地址。

注意： ``mv`` 命令将会删除源文件（对象），如果你不想删除源文件（对象），请
使用 ``mv`` 命令。

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

下面的 ``mv`` 命令将一个本地文件移动到 ``mybucket``::

    $ qsctl mv test.txt qs://mybucket

输出::

    Key <test.txt> created
    File 'test.txt' deleted

下面的 ``mv`` 命令将 ``mybucket`` 下的所有对象移动到本地::

    $ qsctl mv qs://mybucket test/ -r

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.jpg' written
    Key <test/test1.txt> deleted
    Key <test/test2.txt> deleted
    Key <test/test3.jpg> deleted

下面的 ``mv`` 命令将排除所有匹配 ``"*.txt"`` 的对象::

    $ qsctl mv qs://mybucket test -r --exclude "*.txt"

输出::

    File 'test/test3.jpg' written
    Key <test/test3.jpg> deleted

下面的 ``mv`` 命令只包含匹配 ``"*.txt"`` 的对象::

    $ qsctl mv qs://mybucket test -r --exclude "*" --include "*.txt"

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    Key <test/test1.txt> deleted
    Key <test/test2.txt> deleted

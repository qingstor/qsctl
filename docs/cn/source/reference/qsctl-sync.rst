.. _qsctl-sync:


**********
qsctl-sync
**********


====
总览
====

::

      qsctl sync
    <local-path> <qs-path> or <qs-path> <local-path>
    [--delete]
    [--exclude]
    [--include]
    [--rate-limit]

====
描述
====

在本地文件夹和qs-directory间同步。第一个Path变量是源地址，第二个是目标地址。

如果要传的文件（对象）在目标地址中已存在，程序将比较源文件（对象）和目标对
象（文件）的修改时间。只有源文件（对象）比目标对象（文件）新时，才会覆盖目
标对象(文件)。

====
选项
====

``--delete``

加该选项后，目标地址中存在的文件（对象）若源地址中不存在则将被删除。

``--exclude``

排除匹配特定类型的文件（对象）。

``--include``

包含匹配特定类型的文件（对象）。

``--rate-limit``

网速限制,单位可以为: K/M/G，如: 100K 、 1M。

====
示例
====

下面的 ``sync`` 命令将同步本地当前文件夹到 ``mybucket``::

    $ qsctl sync . qs://mybucket

输出::

    Key <test1.txt> created
    Key <test2.txt> created
    Key <test3.txt> created

下面的 ``sync`` 命令将同步 ``QS-Directory`` 到本地::

    $ qsctl sync qs://mybucket/test/ test/

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

下面的 ``sync`` 命令加了 ``--delete`` 选项。假设本地文件夹 ``test/`` 有文件
``test3.txt`` ， ``qs://mybucket/test/`` 下只有 ``test1.txt`` 和
``test2.txt`` , 则本地文件 ``test3.txt`` 将会被删除::

    $ qsctl sync qs://mybucket/test/ test/ --delete

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written
    File 'test/test3.txt' deleted

下面的 ``sync`` 命令将同步本地当前文件夹到 ``mybucket``，并限制速度为每秒 100K::

    $ qsctl sync . qs://mybucket --rate-limit 100K

输出::

    Key <test1.txt> created
    Key <test2.txt> created
    Key <test3.txt> created

下面的 ``sync`` 命令将同步 ``QS-Directory`` 到本地，并限制速度为每秒 100K::

    $ qsctl sync qs://mybucket/test/ test/ --rate-limit 100K

输出::

    File 'test/test1.txt' written
    File 'test/test2.txt' written

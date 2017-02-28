.. _qsctl:


*****
qsctl
*****


====
总览
====

::

    qsctl <command> [<args>] [<options>]

====
选项
====

``-c, --config``

使用特定的配置文件。

``-h, --help``

打印帮助。

``-v``

打印 Qsctl 版本。

====
描述
====

Qsctl是对象存储服务的高级命令行工具。它提供了更强大的类UNIX命令，使管理对象存储 
资源变得像管理本地资源一样方便。这些命令包括：cp，ls，mb，mv，rm，rb和sync。所有
命令都支持批量操作。

下面介绍使用qsctl需要知道的概念。

Path 变量类型
+++++++++++++

Path 变量有两种类型： ``local-path`` 和 ``qs-path`` 。

``local-path``: 本地文件或文件夹的路径。

``qs-path``: QingStor对象, 前缀和存储空间的路径。它的格式写为：
``qs://mybucket/mykey`` ，这里mybucket是存储空间，mykey是QingStor对象。必须
以 ``qs://`` 开头表示这是 ``qs-path`` 。前缀可以在批量操作中使用。比如运行
``qsctl rm -r qs://mybucket/myprefix`` ，所有前缀为 ``myprefix`` 的对象将被
删除。如果前缀为空（如qs://mybucket/），批量操作将对存储空间中所有对象进行。

QS-Directory
++++++++++++
部分命令只对目录进行操作：如 ``sync`` 命令和带 ``-r`` 选项的 ``cp, mv`` 命令。
在qsctl中, 以 '/' 结尾的前缀被视为 ``qs-directory`` 。对于目录操作，如果给定
``qs-path`` 不以 '/' 结束，程序将在末尾自动添加 '/'。即运行
``qsctl sync /foo/bar/ qs://mybucket/myprefix`` 的效果等同于运行
``qsctl sync /foo/bar/ qs://mybucket/myprefix/`` 。空前缀也是
``qs-directory`` （如 ``qs://mybucket`` 和 ``qs://mybucket/`` ），
可以理解成为存储空间的根目录。

Path变量的顺序
++++++++++++++

除了 ``ls`` 之外，所有命令都有一到两个path变量。如果有两个path变量，则第一个
变量代表源地址，第二个代表目的地址。只有一个path变量的命令只对源地址进行操作。

Exclude 和 Include过滤器
++++++++++++++++++++++++

qsctl的通配符操作利用Exclude和Include过滤器实现。它和Unix操作系统的实现有些
区别。比如 ``cp -r dir1/* dir2/`` 这样的操作qsctl不支持。但你可以使用Exclude
和Include过滤器达到同样的效果。Exclude用来排除一些文件（对象），Include用来
包含一些文件（对象）。

下面的通配符是qsctl支持的。

* ``*``: 匹配任意字符串
* ``?``: 匹配任意字符

举例说明, 比如你有这样一个文件结构::

    /tmp/foo/
      bar/
      |---mykey1
      |---mykey2
      test1.txt
      test2.txt
      test3.jpg

运行 ``qsctl cp /tmp/foo qs://mybucket -r --exclude "bar/*"`` ，文件
``bar/mykey1`` 和 ``bar/mykey2`` 将不会被复制到 ``mybucket`` 上。Exclude过滤器
将会检查所有的源文件::

    /tmp/foo/bar/* \-> /tmp/foo/bar/mykey1  （匹配，排除此文件）
    /tmp/foo/bar/* \-> /tmp/foo/bar/mykey2  （匹配，排除此文件）
    /tmp/foo/bar/* \-> /tmp/foo/test1.txt  （不匹配, 包含此文件）
    /tmp/foo/bar/* \-> /tmp/foo/test2.txt  （不匹配, 包含此文件）
    /tmp/foo/bar/* \-> /tmp/foo/test3.jpg  （不匹配, 包含此文件）

注意：默认所有文件是被包含的, 因此单独使用Include过滤器是无效的。Include过滤
器可用来包含被Exclude过滤器过滤掉的文件（对象）。比如如果你只想上传文件夹中
的jpg文件，你可以运行如下命令::

    $ qsctl cp /tmp/foo/ qs://mybucket/ -r --exclude "*" --include "*.jpg"

==========
支持的命令
==========

* ``cp``
* ``ls``
* ``mb``
* ``mv``
* ``rb``
* ``rm``
* ``sync``
* ``presign``

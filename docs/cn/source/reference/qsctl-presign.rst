.. _qsctl-presign:


************************
qsctl-presign
************************


====
总览
====

::

      qsctl presign
    <qs-path>
    [-e, --expire]

====
描述
====

生成 QingStor 对象的参数签名链接。如果是非公开权限，则在请求过期时间前，可以使用该链接下载 QingStor 对象。
如果是公开权限，则在路径不发生改变的情况下，生成的链接永久有效。

====
选项
====

``-e, --expire``

如果是非公开权限，则以-e指定参数签名链接的过期时间，以秒为单位。默认为 3600 秒。仅当是公开权限时，此参数不起作用。

====
示例
====

下面的 ``presign`` 命令生成了 ``mybucket/myobject`` 对象的参数签名链接::

    $ qsctl presign qs://mybucket/myobject

如果是非公开权限，则输出::

    https://pek3a.qingstor.com:443/mybucket/myobject?
    signature=Miy/lgcPTU%2BtzBSnO4nJAdHsEh%2BEo6phvRZc1urckdE%3D
    &access_key_id=EYGSPEMLUGZFBKORUSYO&expires=1488276950

如果是公开权限，则输出::

    https://mybucket.pek3a.qingstor.com/myobject

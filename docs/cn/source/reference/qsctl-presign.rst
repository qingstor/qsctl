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

生成 QingStor 对象的参数签名链接。在请求过期时间前，可以使用该链接下载 QingStor 对象。

====
选项
====

``-e, --expire``

参数签名链接的过期时间，以秒为单位。默认为 3600 秒。

====
示例
====

下面的 ``presign`` 命令生成了 ``mybucket/myobject`` 对象的参数签名链接::

    $ qsctl presign qs://mybucket/myobject

输出::

    https://pek3a.qingstor.com:443/mybucket/myobject?
    signature=Miy/lgcPTU%2BtzBSnO4nJAdHsEh%2BEo6phvRZc1urckdE%3D
    &access_key_id=EYGSPEMLUGZFBKORUSYO&expires=1488276950


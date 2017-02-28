.. _qsctl-presign:


************************
qsctl-presign
************************


========
Synopsis
========

::

      qsctl presign
    <qs-path>
    [-e, --expire]

===========
Description
===========

Generate a pre-signed URL for an object. Within the given expire time,
anyone who receives this URL can retrieve the object with an HTTP GET request.

=======
Options
=======

``-e, --expire``

The number of seconds until the pre-signed URL expires. Default is 3600 seconds.

========
Examples
========

The following ``presign`` command presigns object ``mybucket/myobject``::

    $ qsctl presign qs://mybucket/myobject

Output::

    https://pek3a.qingstor.com:443/mybucket/myobject?
    signature=Miy/lgcPTU%2BtzBSnO4nJAdHsEh%2BEo6phvRZc1urckdE%3D
    &access_key_id=EYGSPEMLUGZFBKORUSYO&expires=1488276950


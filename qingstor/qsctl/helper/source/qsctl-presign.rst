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

If an object belongs to a non-public bucket, generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve the object with
an HTTP GET request. If an object belongs to a public bucket, generate a URL spliced by
bucket name, zone and its name, anyone who receives this URL can always retrieve the
object with an HTTP GET request.

=======
Options
=======

``-e, --expire``

If the object belongs to a non-public bucket, the parameter means the number of seconds until
the pre-signed URL expires. Default is 3600 seconds.If the object belongs to a public bucket,
the parameter has no effects.

========
Examples
========

The following ``presign`` command presigns object ``mybucket/myobject``::

    $ qsctl presign qs://mybucket/myobject

If the object belongs to a non-public bucket, then output::

    https://pek3a.qingstor.com:443/mybucket/myobject?
    signature=Miy/lgcPTU%2BtzBSnO4nJAdHsEh%2BEo6phvRZc1urckdE%3D
    &access_key_id=EYGSPEMLUGZFBKORUSYO&expires=1488276950

If the object belongs to a public bucket, then output::

    https://mybucket.pek3a.qingstor.com/myobject

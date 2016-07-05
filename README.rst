=====
qsctl
=====

qsctl is intended to be an advanced command line tool for QingStor, it provides
powerful unix-like commands to let you manage QingStor resources just like files
on local machine. Unix-like commands contains: cp, ls, mb, mv, rm, rb, and sync.
All of them support batch processing.

------------
Installation
------------

virtualenv::

    $ pip install qsctl

System-Wide::

    $ sudo pip install qsctl

On Windows systems, run it in a command-prompt window with administrator
privileges, and leave out sudo.

---------------
Getting Started
---------------

To use qsctl, there must be a configuration file to configure your own
``qy_access_key_id``, ``qy_secret_access_key``, and ``zone``, for example::

  qy_access_key_id: 'QINGCLOUDACCESSKEYID'
  qy_secret_access_key: 'QINGCLOUDSECRETACCESSKEYEXAMPLE'
  zone: 'pek3a'

The configuration file is ``~/.qingcloud/config.yaml`` by default, it also
can be specified by the option ``-c /path/to/config``.

------------------
Available Commands
------------------

Commands supported by qsctl are listed below:

.. list-table::
  :widths: 10 90
  :header-rows: 0

  * - ls
    - List qingstor keys under a prefix or all qingstor buckets.

  * - cp
    - Copy local file(s) to qingstor or qingstor key(s) to local.

  * - mb
    - Create a qingstor bucket.

  * - rb
    - Delete an empty qingstor bucket or forcibly delete nonempty qingstor bucket.

  * - mv
    - Move local file(s) to qingstor or qingstor keys(s) to local.

  * - rm
    - Delete a qingstor key or keys under a prefix.

  * - sync
    - Sync between local directory and qingstor prefix.

--------
Examples
--------

List keys in bucket <mybucket> by running::

  $ qsctl ls qs://mybucket
  Directory                          test/
  2016-04-03 11:16:04     4 Bytes    test1.txt
  2016-04-03 11:16:04     4 Bytes    test2.txt

Sync from qingstor prefix to local directory::

  $ qsctl sync qs://mybucket3/test/ test/
  File 'test/README.md' written
  File 'test/commands.py' written

See the detailed usage and more examples with 'qsctl help' or 'qsctl <command> help'.

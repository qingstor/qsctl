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

To use qsctl, there must be a configuration file , for example::

  access_key_id: 'ACCESS_KEY_ID_EXAMPLE'
  secret_access_key: 'SECRET_ACCESS_KEY_EXAMPLE'

The configuration file is ``~/.qingstor/config.yaml`` by default, it also
can be specified by the option ``-c /path/to/config``.

You can also config other option like ``host`` , ``port`` and so on, just
add lines below into configuration file, for example::

  host: 'qingstor.com'
  port: 443
  protocol: 'https'
  connection_retries: 3
  # Valid levels are 'debug', 'info', 'warn', 'error', and 'fatal'.
  log_level: 'debug'

------------------
Available Commands
------------------

Commands supported by qsctl are listed below:

.. list-table::
  :widths: 10 90
  :header-rows: 0

  * - ls
    - List QingStor keys under a prefix or all QingStor buckets.

  * - cp
    - Copy local file(s) to QingStor or QingStor key(s) to local.

  * - mb
    - Create a QingStor bucket.

  * - rb
    - Delete an empty QingStor bucket or forcibly delete nonempty QingStor bucket.

  * - mv
    - Move local file(s) to QingStor or QingStor key(s) to local.

  * - rm
    - Delete a QingStor key or keys under a prefix.

  * - sync
    - Sync between local directory and QingStor prefix.

  * - presign
    - Generate a pre-signed URL for an object.

--------
Examples
--------

List keys in bucket <mybucket> by running::

  $ qsctl ls qs://mybucket
  Directory                          test/
  2016-04-03 11:16:04     4 Bytes    test1.txt
  2016-04-03 11:16:04     4 Bytes    test2.txt

Sync from QingStor prefix to local directory::

  $ qsctl sync qs://mybucket3/test/ test/
  File 'test/README.md' written
  File 'test/commands.py' written

See the detailed usage and more examples with 'qsctl help' or 'qsctl <command> help'.

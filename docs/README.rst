==========================
Building The Documentation
==========================

Before building the documentation, make sure you have Python 2.7, Qsctl, and
Sphinx installed. Sphinx can be installed by::

    pip install sphinx

We have both Chinese and English documentation. The configuration files are
in ``cn`` and ``en`` directories. Please build the documentation in the
corresponding directory.

The methods to build the documentation:

* ``make html`` to build the documentation in HTML format into the
  ``build/html`` directory.

* ``make man`` to build the documentation in Linux man page format into
  ``../doc/man``.

* ``make text`` to build the documentation in text format that can be used
  on Windows platform.

You can also build the documentation in other formats supported by Sphinx.
See `Sphinx docs <http://zh-sphinx-doc.readthedocs.io/en/latest/contents.html>`_.

===========
 Slackdump
===========

.. contents::

Installation
------------

Installing is pretty simple - just download the latest Slackdump from the
Releases page, extract and run it:

#. Download the archive from the Releases_ page for your operating system.

   .. tip:: **MacOS users** should download ``darwin`` release file.
#. Unpack;
#. Change directory to where you have unpacked the archive;
#. Run ``./slackdump -h`` to view help options.


Logging in
----------
There are two types of login options available:

- Automatic_ (also called **EZ-Login 3000**); OR
- Manual_

Automatic_ login is the default one, it requires no prior setup, and the
general recommendation is to use the Automatic login.  If the Automatic login
doesn't work for some reason, fallback to Manual_ login steps.

Usage
-----
There are three modes of operation:

- `Listing users/channels`_
- `Dumping messages and threads`_ (private and public)
- `Creating a Slack Export`_


.. _Automatic:  login-auto.rst
.. _Manual: login-manual.rst
.. _Installation: usage-install.rst
.. _Dumping messages and threads: usage-channels.rst
.. _Creating a Slack Export: usage-export.rst
.. _Listing users/channels:  usage-list.rst
.. _Releases: https://github.com/rusq/slackdump/releases

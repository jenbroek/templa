templa - Go templating utility
==============================
templa reads from standard input, executing it as a Go template
using values from the specified value files (if any).


Value files
-----------
Currently only YAML is supported.


Installation
------------
Edit config.mk to match your local setup (templa is installed into
the /usr/local namespace by default).

Afterwards enter the following command to build and install templa (if
necessary as root):

	make clean install


Running templa
--------------
Simply invoke 'templa', optionally redirecting input/output.
See the man page for details.

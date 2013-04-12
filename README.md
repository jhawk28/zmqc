zmqc
====

my go port of the https://github.com/zacharyvoase/zmqc

I wanted a simple cat style utility that was easier to install (no runtime dependencies other than zeromq). It mostly functions the same as the python zmqc, but I simplified the flags and do not support the generic socket options.

Tested compile on Windows and OSX.

Depends on:
github.com/droundy/goopt
github.com/alecthomas/gozmq

To install:
go get github.com/jhawk28/zmqc

Options:
  -0                                                 Separate messages on input/output should be 
  -c                                   --connect     Connect to the specified address(es).
  -b                                   --bind        Bind to the specified address(es).
  -n -1                                              Receive/send only NUM messages. By default, zmqc lives forever in 'read' mode, or until the end of input in 'write' mode.
  -s [PUSH|PULL|PUB|SUB|REQ|REP|PAIR]                Which type of socket to create. Must be one of 'PUSH', 'PULL', 'PUB', 'SUB', 'REQ', 'REP' or 'PAIR'. See `man zmq_socket` for an explanation of the different types. 'DEALER' and 'ROUTER' sockets are currently unsupported.
                                       --subscribe=  Subscribes to data matching
                                       --help        show usage message


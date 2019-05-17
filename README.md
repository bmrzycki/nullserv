# nullserv

## What is it?
It's a simple null http and https server originally written using Go 1.5.
The language has evolved considerably since then but I've tried to keep
the code as close to this version as I can. This also means I haven't
pulled in any third-party packages to minimize dependency issues on
older embedded devices.

## Why would I want it?
Because you're running a DNS ad blocker and you want a server that
understands several file extensions and returns cached, minimal files for
each.

## How do I install it?
Pull the repo, install Google Go and run make. Why use make instead of
go build? I need to dynamically generate version.go as well as compile
a small helper program (written in clean ANSI C) that generates Go's
[]byte{...} syntax.

If the idea of using make is abhorrent to you then you can just run:
```
$ ./go_ver.sh
$ go build -o nullsrv *.go
```

## Lower numbered-ports
Listening on TCP ports lower than 1024 usually requires special OS access.
On Linux, you can either run as root or use setcap to give the created
nullserv binary access to low number ports. It's also a good idea to use
a daemon launcher like start-stop-daemon to run it.

On my setup I perform the following actions to deploy nullserv and run it
as user nobody.

1. Compile nullserv: $ make

2. Copy the binary as root: # cp nullserv /usr/local/bin

3. Change the permissions: chown root:nogroup /usr/local/bin/nullserv; chmod 750 /usr/local/bin/nullserv

4. Use Linux capabilities: setcap 'cap_net_bind_service=+ep' /usr/local/bin/nullserv

5. Launch the binary: /sbin/start-stop-daemon -S -b -c nobody:nogroup -x /usr/local/bin/nullserv

## Command line interface
```
$ ./nullserv -h
Usage of ./nullserv:
  -A string
        https address (default all)
  -P int
        https port (default 443)
  -a string
        http address (default all)
  -m int
        content cache age in secs (default 604800)
  -p int
        http port (default 80)
  -v int
        verbose 0..9 (default 0)
```

# Background
The idea for nullserv came from another GitHub project,
[pixelserv](https://github.com/h0tw1r3/pixelserv). I wanted a real project
to learn how to use Google Go. The running binary on a 64-bit x86 host
consumes about 8MB of real memory, mostly from Go libraries. It's not as
small as pixelserv is, but I think it's a lot easier to extend and
maintain. Go's low-level http library and go routines make scalability
a breeze.

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

When listening to the http port nullserv returns small, but valid,
responses for several common file suffixes. This reduces web page layout
problems. For mobile apps a valid file is often intepreted as a "good"
ad and will effectively remove annoying ads from several mobile apps as
well as traditional browsers.

When litening to the https port, nullserv aborts the connection early
in a fairly graceful manner. This side-steps the requirement for
proxies or self-generated certificate authorities. It's not as elegant
as http traffic responses but it doesn't require changes to the security
of client browsers or OSes.

## How do I install it?
Pull the repo, install [Google Go](https://golang.org/) and run
`make`. Why use `make` instead of `go build`? I need to dynamically
generate `version.go` as well as compile a small helper program written in
clean ANSI C that generates Go's `[]byte{...}` array syntax similar to what
`xxd -i filename` emits for C headers.

If the idea of using `make` is abhorrent to you then you can just run:
```
$ ./go_ver.sh
$ go build -o nullserv *.go
```

## Lower numbered-ports
Listening on TCP ports lower than 1024 usually requires special OS access.
On Linux, you can either run as root or use setcap to give the created
nullserv binary access to low number ports. It's also a good idea to use
a daemon launcher like start-stop-daemon to run it.

On my setup I perform the following actions to deploy nullserv and run it
with reduced permissions.

```
# Build
$ make clean && make -j

# Copy binary and change the binary permissions
$ sudo cp nullserv /usr/local/bin
$ sudo chown root:nogroup /usr/local/bin/nullserv
$ sudo chmod 750 /usr/local/bin/nullserv

# Use Linux capabilities
$ sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/nullserv

# Launch
$ sudo /sbin/start-stop-daemon -S -b -c nobody:nogroup -x /usr/local/bin/nullserv
```

## Command line interface
```
Usage of ./nullserv:
  -A string
    	https address (default '' = all)
  -P int
    	https port (default 443)
  -a string
    	http address (default '' = all)
  -c string
    	JSON config file
  -m int
    	content cache age in secs (default 31536000)
  -p int
    	http port (default 80)
  -v int
    	verbose level
```

## Config files
You can use JSON files as a configuration file for nullserv. There are several
examples in the example_confs/ subdirectory.  Note that the contents of a
configuration file overrides what is passed in on the command line. Here's
a simple config file when listening on the standard ports on all interfaces:
```
{
    "max_age" : 31536000,
    "verbose" : 0,
    "http" : {
	"address" : "",
	"port" : 80
    },
    "https" : {
	"address" : "",
	"port" : 443
    }
}
```

# Background
The idea for nullserv came from another GitHub project,
[pixelserv](https://github.com/h0tw1r3/pixelserv). I wanted a real project
to learn how to use Google Go. The running binary on a 64-bit x86 host
consumes about 8MB of real memory, mostly from Go libraries. It's not as
small as pixelserv is, but I think it's a lot easier to extend and
maintain. Go's low-level http library and go routines make scalability
a breeze.

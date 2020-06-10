# nullserv

A simple null file http and https server originally written using Go 1.5.

Go has evolved considerably since 1.5 but I've tried to keep
the code as close to this version. I also haven't used any third-party
packages to minimize dependency build issues.  If your platform has
Go then `nullserv` will likely run on it.

## Why would I want it?
Because you're running a DNS ad blocker and you want a server that
returns minimal valid files for each based on file suffixes in the URL.

`nullserv` reduces web page layout problems when blocking ads. It also
fixes layout and "good" ad detection for moble apps because the responses
from `nullserv` are valid images.

The modern web as moved on since I started this project with almost all
websites using encrypted https traffic for ads and trackers. The good
news is `nullserv` also listens on https port 443 and aborts all SSL/TLS
connections early. It claims the client's SSL certificate of authority (CA)
is invalid. This side-steps the requirement for proxies or self-generated
certificate authorities needed to run a real https server.

## How do I install it?
Pull the repo, install [Google Go](https://golang.org/) and run
`make`. I use `make` instead of `go build` to dynamically
generate `buildinfo.go` as well as `file2gobyte`.

If you prefer to manually build `nullserv` run:
```
$ cd /path/to/cloned/repo/nullserv
$ cc -Wall -O3 -o file2gobyte file2gobyte.c  # OPTIONAL
$ sh mkbuildinfo.sh
$ go build -o nullserv *.go
```

The optional `file2gobyte` utility emits Go's `[]byte{...}` array syntax
on stdout similar to `xxd -i filename` for C. It's not necessary to
run `nullserv` and only used to assist adding file data to `files.go`.

There's also a `Dockerfile` to build and run `nullserv` inside a Docker
instance.

The emitted `nullserv` binary is self-contained and relies on no external
files (except maybe a JSON config file). You may run it from within the
repo or copy it to somewere common like `/usr/local/bin`.

## Running on low numbered TCP ports as non-root
It's recommended to run daemons as a non-root user whenever possible.
Unfortunately `nullserv` bind to TCP ports 80 (http) and 443 (https)
by default. Most Unix operating systems forbid anyone but root to bind to
privileged ports lower than 1024. The [Go language also has
problems](https://github.com/golang/go/issues/1435) with `Setuid/Setgid`
across all threads which prevents handling the issue inside `nullserv`.

The simplest way to run `nullserv` as non-root is to temporarily disable
privileged ports checking at the OS level, run `nullserv` as a daemon,
and re-enable privileged port checking. Here's a Linux eample:

```
# (as user root or via sudo)
# Disable unprivileged port binding on 80 and 443 to allow nullserv
# to run as user nobody. After the daemon starts re-enable on < 1024.
# Even though this is in the net.ipv4 namespace it also allows binding
# of IPv6 ports.
sysctl -q net.ipv4.ip_unprivileged_port_start=80

# start-stop-daemon if found in any Debian, or openrc, based distro.
# Other distros have comperable dameonizing tools. Red Hat based
# distros can use a combination of 'nohup' and 'runner'.
start-stop-daemon --start --background \
    --exec "/path/to/binary/nullserv" \
    --stdout "$LOG/nullserv.log" --stderr "$LOG/nullserv.log" \
    --user nobody

# re-enable low port checking
sysctl -q net.ipv4.ip_unprivileged_port_start=1024
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
`nullserv` supports JSON configuration files as an alternative to command
line arguments. See the `example_confs/` subdirectory for a few use-cases.
The contents of the file takes precedence over command line arguments and
any missing parameters use `nullserv` defaults.

For example if you only wish to change the `max_age` parameter to use
`no-store` and never have clients cache responses:

```
$ cat max_age.conf
{ "max_age" : -1 }
$ nullserv -c max_age.conf
```

The `ConfFile` struct in `conf.go` lists all the valid JSON parameters
names as does the JSON file `example_confs/all.json`.

## Internal running state
`nullserv` interprets certain file suffixes as requests for internal state.

### Version and build info (.ver)
Any URL ending with `.ver` or `.version` will return version information in
JSON:

```
$ curl http://127.0.0.1/.ver
{
  "build_date": "Wed, 10 Jun 2020 09:44:58 -0500",
  "commit_date": "Sat, 30 May 2020 17:33:56 -0500",
  "reset_date": "Wed, 10 Jun 2020 12:45:34 -0500",
  "sha": "56596c6",
  "version": "1.3.0"
}
```

### Statistics (.stat)
Any URL ending with `.stat` or `.stats` will return statistics:

```
$ curl http://127.0.0.1/.stat
{
  "http_1.1": 6,
  "https_tls_1.0": 9,
  "stats_ok": 1,
  "stats_version": 3,
  "suffix_" : 1,
  "suffix_stats": 3,
  "suffix_version": 2
}
```

Keys starting with `http_` and `https_` track the client protocol requests
counts. In the example above we've seen 6 http 1.1 and 9 https 1.0
connections.

The `stats_` keys are meta-data. The `stats_ok` value returns `1` if
the `.stats` structure is valid, correct data (`0` otherwise). The
`stats_version` is a monotonically increasing version number which will
be incremented whenever the layout of this JSON response is altered.

The `suffix_*` keys show counts for each of the file suffixes encountered
since the last `.reset` event. The above data has 1 request for a URL
with no file suffix, 3 for file suffix `.stats`, and 2 for file
suffix `.version`.

### Resetting statistics (.res)
Any URL ending with `.res` or `.reset` will reset all internal statistics.
The response is identical to `.stat` after resetting all counts. This also
updates `reset_date` in the `.ver` output to the current time on the host
running `nullserv`.

```
$ curl http://127.0.0.1/.res
{
  "stats_ok": 1,
  "stats_version": 3
}
```

## Background
The idea for nullserv came from another GitHub project,
[pixelserv](https://github.com/h0tw1r3/pixelserv). I wanted a real project
to learn how to use Google Go. The running binary on a 64-bit x86 host
consumes about 8MB of real memory (mostly Go libraries). It's not as
small as pixelserv is but I think it's a lot easier to extend and
maintain. Go's low-level `http` library and go routines make scalability
a breeze.

package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type nullCounter map[string]int

type logMsg struct {
	pri syslog.Priority
	str string
}

const pkgName = "nullserv"

var maxAge int
var verbose int
var stats chan string
var msg chan logMsg

func doLog(p syslog.Priority, s string) {
	msg <- logMsg{
		pri: p,
		str: s,
	}
}

// Inspired by:
//   https://github.com/h0tw1r3/pixelserv/blob/master/pixelserv.c
func fakeHTTPS(l net.Listener) (err error) {
	conn, err := l.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	if verbose > 0 {
		doLog(syslog.LOG_NOTICE, "HTTPS (?)")
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	// Always respond with a TLS access denied error
	_, err = conn.Write([]byte{
		'\x15',         // Alert 21
		'\x03', '\x00', // Version 3.0
		'\x00', '\x02', // Length == 2
		'\x02', // Fatal event
		'\x31', // 0x31 == TLS access denied (49)
	})
	if err != nil {
		return err
	}
	stats <- "(https)"
	return nil
}

func adservHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.QueryUnescape(r.URL.String())

	if verbose > 0 {
		doLog(syslog.LOG_NOTICE,
			fmt.Sprintf("HTTP %s %s", r.Host, u))
	}

	for _, value := range []string{"?", ";", "#"} {
		if strings.Contains(u, value) {
			u = strings.Split(u, value)[0]
		}
	}

	suffix := ""
	if strings.Contains(u, ".") {
		tmp := strings.Split(u, ".")
		suffix = strings.ToLower(tmp[len(tmp)-1])
	}

	if verbose > 0 {
		doLog(syslog.LOG_NOTICE,
			fmt.Sprintf("HTTP (%s)", suffix))
	}

	if _, ok := NotFoundFiles[suffix]; ok == true {
		http.NotFound(w, r)
		stats <- "(not found)"
		return
	}

	// Locate alternate suffix spellings or related file types
	if realSuffix, ok := AltSuffix[suffix]; ok == true {
		suffix = realSuffix
	}

	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(maxAge))
	if f, ok := nullFiles[suffix]; ok == true {
		w.Header().Set("Content-Type", f.content)
		if f.data != nil {
			w.Write(f.data)
		}
	} else {
		f = nullFiles["html"] // unknown suffixes become HTML
		w.Header().Set("Content-Type", f.content)
		w.Write(f.data)
	}
	stats <- suffix
}

func statHandler() {
	count := make(map[string]int)
	for {
		suffix := <-stats
		if suffix == "__show__" {
			doLog(syslog.LOG_NOTICE, fmt.Sprintf("%v", count))
			count = make(map[string]int) // reset all counters
		} else {
			if _, ok := count[suffix]; ok == true {
				count[suffix]++
			} else {
				count[suffix] = 1
			}
		}
	}
}

func logHandler() {
	useLog := false
	sl, err := syslog.New(syslog.LOG_NOTICE, pkgName)
	if err != nil {
		useLog = true
	}
	defer sl.Close()
	for {
		m := <-msg
		switch {
		case useLog == true:
			log.Println(m.str)
		case m.pri == syslog.LOG_ERR:
			sl.Err(m.str)
		default:
			sl.Notice(m.str)
		}
	}
}

func main() {
	// Setup command line arguments
	httpAddrPtr := flag.String("a", "", "http address (default all)")
	httpPortPtr := flag.Int("p", 80, "http port")
	httpsAddrPtr := flag.String("A", "", "https address (default all)")
	httpsPortPtr := flag.Int("P", 443, "https port")
	maxAgePtr := flag.Int("m", 604800, "content cache age in secs")
	verbosePtr := flag.Int("v", 0, "verbose 0..9 (default 0)")
	flag.Parse()

	maxAge = *maxAgePtr
	verbose = *verbosePtr

	stats = make(chan string, 10) // arbitrary size: grow when prog pauses
	go statHandler()

	msg = make(chan logMsg, 10) // arbitrary size: grow when prog pauses
	go logHandler()

	// Starting HTTP server
	addr := *httpAddrPtr + ":" + strconv.Itoa(*httpPortPtr)
	http.HandleFunc("/", adservHandler)
	go func() {
		doLog(syslog.LOG_NOTICE, "Starting HTTP service on "+addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			doLog(syslog.LOG_ERR,
				"http.ListenAndServe error "+err.Error())
		}
	}()

	// Starting fake HTTPS server
	sslAddr := *httpsAddrPtr + ":" + strconv.Itoa(*httpsPortPtr)
	doLog(syslog.LOG_NOTICE, "Starting fake HTTPS service on "+sslAddr)
	l, err := net.Listen("tcp", sslAddr)
	if err != nil {
		doLog(syslog.LOG_ERR,
			"net.Listen HTTPS error "+err.Error())
	}
	go func() {
		for {
			if err := fakeHTTPS(l); err != nil {
				doLog(syslog.LOG_ERR, "fakeHTTPS error "+err.Error())
			}
		}
	}()

	// Starting signal listener
	done := make(chan bool, 1)
	go func() {
		sigChan := make(chan os.Signal)
		for {
			signal.Notify(sigChan, syscall.SIGTERM,
				syscall.SIGUSR1, syscall.SIGUSR2)
			sig := <-sigChan
			if sig == syscall.SIGTERM {
				doLog(syslog.LOG_NOTICE, "Exiting on SIGTERM")
				done <- true
			} else if sig == syscall.SIGUSR1 {
				stats <- "__show__"
			} else if sig == syscall.SIGUSR2 {
				verbose++
				if verbose > 9 {
					verbose = 0
				}
				doLog(syslog.LOG_NOTICE,
					"debug level "+strconv.Itoa(verbose))
			}
		}
	}()
	<-done
}

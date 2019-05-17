package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type SafeCounter struct {
	v   map[string]uint64
	mux sync.Mutex
}

var MaxAgeVal string
var Stats SafeCounter

func AbortTLSListener(conn net.Conn) {
	// This sends a TLS v1.2 alert packet regardless of query.
	// We respond the certificate authority (CA) that issued
	// the requesters certificate is unknown to us. This is possibly
	// the shortest response to a TLS (HTTPS) connection which
	// initiates a graceful shutdown on both ends.
	//
	// Note: This isn't TLS 1.3 compatible. Once browsers start
	//       deprecating 1.2 this will need to be overhauled. The
	//       1.3 specification is nothing like 1.2.
	//
	// The original idea came from h0tw1r3. Relevant projects:
	//  https://github.com/kvic-z/pixelserv-tls/wiki/Command-Line-Options
	//  https://github.com/h0tw1r3/pixelserv/blob/master/pixelserv.c
	conn.Write([]byte{
		'\x15',         // Alert protocol header (21)
		'\x03', '\x03', // TLS v1.2 (RFC 5246)
		'\x00', '\x02', // Message length (2)
		'\x02',  // Alert level fatal (2)
		'\x30'}) // Unknown Certificate Authority (48)
	conn.Close()

	Stats.mux.Lock()
	defer Stats.mux.Unlock()
	Stats.v["_transport_https"]++
}

func NullHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.QueryUnescape(r.URL.String())

	// RFC 3986, Section 3 lists '?' as a query delimiter,
	// '#' as a fragment delimiter, and ';' as a sub-delimiter.
	// All three cannot be part of the path. Remove them.
	for _, value := range []string{"?", ";", "#"} {
		if strings.Contains(u, value) {
			u = strings.Split(u, value)[0]
		}
	}

	// Obtain the file suffix in the URI, if any.
	suffix := ""
	if idx := strings.LastIndex(u, "."); idx != -1 {
		suffix = u[idx+1 : len(u)]
	}

	// If this is an alternate suffix, replace with the real one.
	if realSuffix, ok := AltSuffix[suffix]; ok == true {
		suffix = realSuffix
	}

	Stats.mux.Lock()
	defer Stats.mux.Unlock()
	Stats.v["_transport_http"]++
	Stats.v[suffix]++

	// Special suffix "stats" emits stats about nullsrv as JSON.
	if suffix == "stats" {
		w.Header().Set("Cache-Control", "max-age=0")
		w.Header().Set("Content-Type", "application/json")
		json, err := json.MarshalIndent(Stats.v, "", "    ")
		if err != nil {
			w.Write([]byte("{}"))
		} else {
			w.Write(json)
		}
		return
	}

	// These are suffixes where we return 404, not found.
	if _, ok := NotFoundFiles[suffix]; ok == true {
		http.NotFound(w, r)
		return
	}

	// Fetch the null file for this suffix. Use HTML as the default case.
	f, ok := NullFiles[suffix]
	if ok != true {
		f = NullFiles["html"]
	}
	w.Header().Set("Cache-Control", MaxAgeVal)
	w.Header().Set("Content-Type", f.content)
	if f.data != nil {
		w.Write(f.data)
	}
}

func main() {
	// Parse command line arguments
	httpAddr := flag.String("a", "", "http address (default '' = all)")
	httpPort := flag.Int("p", 80, "http port")
	httpsAddr := flag.String("A", "", "https address (default '' = all)")
	httpsPort := flag.Int("P", 443, "https port")
	maxAge := flag.Int("m", 31536000, "content cache age in secs")
	flag.Parse()

	// Initialize globals
	MaxAgeVal = "public, max-age=" + strconv.Itoa(*maxAge)
	Stats = SafeCounter{v: make(map[string]uint64)}

	// Starting HTTP server
	addr := *httpAddr + ":" + strconv.Itoa(*httpPort)
	http.HandleFunc("/", NullHandler)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal("HTTP service error: " + err.Error())
		}
	}()

	// Starting the abort TLS (HTTPS) server
	sslAddr := *httpsAddr + ":" + strconv.Itoa(*httpsPort)
	l, err := net.Listen("tcp", sslAddr)
	if err != nil {
		log.Fatal("Abort TLS listen error: " + err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Abort TLS accept error: " + err.Error())
		}
		go AbortTLSListener(conn)
	}
}

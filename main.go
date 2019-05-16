package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var maxAge string

func abortTLS(conn net.Conn) {
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
	_, _ = conn.Write([]byte{
		'\x15',         // Alert protocol header (21)
		'\x03', '\x03', // TLS v1.2 (RFC 5246)
		'\x00', '\x02', // Message length (2)
		'\x02',         // Alert level fatal (2)
		'\x30'})        // Unknown Certificate Authority (48)
	conn.Close()
}

func nullHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.QueryUnescape(r.URL.String())

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

	// Locate alternate suffix spellings or related file types
	if realSuffix, ok := AltSuffix[suffix]; ok == true {
		suffix = realSuffix
	}

	if _, ok := NotFoundFiles[suffix]; ok == true {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=" + maxAge)
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
}

func main() {
	// Parse command line arguments
	httpAddr := flag.String("a", "", "http address (default '' = all)")
	httpPort := flag.Int("p", 80, "http port")
	httpsAddr := flag.String("A", "", "https address (default '' = all)")
	httpsPort := flag.Int("P", 443, "https port")
	tmpMaxAge := flag.Int("m", 31536000, "content cache age in secs")
	flag.Parse()
	maxAge = strconv.Itoa(*tmpMaxAge)

	// Starting HTTP server
	addr := *httpAddr + ":" + strconv.Itoa(*httpPort)
	http.HandleFunc("/", nullHandler)
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
		go abortTLS(conn)
	}
}

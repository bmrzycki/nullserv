package main

import (
	"encoding/json"
	"fmt"
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

const TLSbufLen = 3
const TLStypeHandshake = 22
const TLSversionMajor = 3

func AbortTLSListener(conn net.Conn) {
	transport := "https_invalid"
	tls := false
	defer conn.Close()

	// The first 3 bytes determine a properly-formed TLS handshake.
	//   byte 0  ContentType)     SSL/TLS record type
	//   byte 1  ProtocolVersion  SSL/TLS version major
	//   byte 2  ProtocolVersion  SSL/TLS version minor
	// https://tools.ietf.org/html/rfc5246#appendix-A.1
	// https://serializethoughts.com/2014/07/27/dissecting-tls-client-hello-message
	buf := make([]byte, TLSbufLen)
	n, err := conn.Read(buf)
	if err == nil && n == TLSbufLen &&
		buf[0] == TLStypeHandshake && buf[1] == TLSversionMajor {
		switch buf[2] {
		case 0:
			transport = "https_ssl_3.0"
		case 1:
			transport = "https_tls_1.0"
			tls = true
		case 2:
			transport = "https_tls_1.1"
			tls = true
		case 3:
			transport = "https_tls_1.2"
			tls = true
		case 4:
			transport = "https_tls_1.3"
			tls = true
		}
	}

	Stats.mux.Lock()
	defer Stats.mux.Unlock()
	Stats.v[transport]++

	if tls == false {
		return
	}

	// Send a TLS v1.0 alert response packet. Respond the certificate
	// authority (CA) issuer of the client certificate is unknown to us.
	// This is possibly the shortest response to a TLS connection which
	// gracefully terminates a TLS handshake.
	//
	// TLS is intentionally backward compatible to 1.0 and is why we
	// can always send it regardless of the client request version.
	// https://blog.cloudflare.com/why-tls-1-3-isnt-in-browsers-yet/
	//
	// The original idea came from h0tw1r3's pixelserv:
	//  https://github.com/h0tw1r3/pixelserv/blob/master/pixelserv.c
	conn.Write([]byte{
		'\x15',         // Alert protocol header (21)
		'\x03', '\x01', // TLS v1.0 (RFC 5246)
		'\x00', '\x02', // Message length (2)
		'\x02',  // Alert level fatal (2)
		'\x30'}) // Unknown Certificate Authority (48)
}

func NullHandler(w http.ResponseWriter, r *http.Request) {
	Stats.mux.Lock()
	defer Stats.mux.Unlock()

	u, err := url.QueryUnescape(r.URL.String())
	if err != nil {
		Stats.v["http_invalid_url"]++
		return
	}

	// RFC 3986, Section 3 lists '?' as a query delimiter,
	// '#' as a fragment delimiter, and ';' as a sub-delimiter.
	// All three must be stripped from the url.
	for _, value := range []string{"?", ";", "#"} {
		if strings.Contains(u, value) {
			u = strings.Split(u, value)[0]
		}
	}

	// Obtain the file suffix in the URI, if any.
	suffix := ""
	if idx := strings.LastIndex(u, "."); idx != -1 {
		suffix = u[idx+1 : len(u)]
		if real, ok := AltSuffix[suffix]; ok == true {
			suffix = real
		}
	}

	// Special suffix ".reset" resets statistics.
	if suffix == "reset" {
		for k := range Stats.v {
			delete(Stats.v, k)
		}
		suffix = "stats"
		GenVersion()
	} else {
		s := fmt.Sprintf("http_%d.%d", r.ProtoMajor, r.ProtoMinor)
		Stats.v[s]++
		Stats.v["suffix_"+suffix]++
	}

	// Handle the 404 not found suffixes.
	if _, ok := NotFoundFiles[suffix]; ok == true {
		http.NotFound(w, r)
		return
	}

	cc := MaxAgeVal
	if suffix == "version" {
		cc = "no-store"
	}

	// Obtain data with HTML as default.
	f, ok := NullFiles[suffix]
	if ok != true {
		f = NullFiles["html"]
	}
	data := f.data

	// Generate new json stats if requested.
	if suffix == "stats" {
		cc = "no-store"
		Stats.v["ok"] = 1
		json, err := json.MarshalIndent(Stats.v, "", "  ")
		if err == nil {
			data = json
		}
	}

	w.Header().Set("Cache-Control", cc)
	w.Header().Set("Content-Type", f.content)
	if data != nil {
		w.Write(data)
	}
}

func main() {
	// Initialize globals
	ConfInit()
	Stats = SafeCounter{v: make(map[string]uint64)}
	GenVersion()
	if Config.MaxAge == -1 {
		MaxAgeVal = "no-store"
	} else {
		MaxAgeVal = "public, max-age=" + strconv.Itoa(Config.MaxAge)
	}

	// Starting HTTP server
	a := Config.Http.Address + ":" + strconv.Itoa(Config.Http.Port)
	http.HandleFunc("/", NullHandler)
	go func() {
		if err := http.ListenAndServe(a, nil); err != nil {
			log.Fatal("HTTP service error: " + err.Error())
		}
	}()

	// Starting the abort TLS (HTTPS) server
	sa := Config.Https.Address + ":" + strconv.Itoa(Config.Https.Port)
	l, err := net.Listen("tcp", sa)
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

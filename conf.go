package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type ConfFile struct {
	MaxAge  int `json:"max_age"`
	Verbose int `json:"verbose"`
	Http    struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}
	Https struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}
}

const HTTPMaxAge = 31536000 // RFC 2616, Section 14.21
const TCPMaxPort = 65535    // RFC 793, Sction 3.1, 16-bit port size

var Config ConfFile

func ConfValidate() {
	if Config.MaxAge < 0 || Config.MaxAge > HTTPMaxAge {
		log.Fatalf("MaxAge must be in range 0..%d", HTTPMaxAge)
	}

	if Config.Verbose < 0 {
		log.Fatal("Verbose must be 0 or larger")
	}

	if Config.Http.Port < 0 || Config.Http.Port > TCPMaxPort {
		log.Fatalf("HTTP port must be in range 0..%d", TCPMaxPort)
	}

	if Config.Https.Port < 0 || Config.Https.Port > TCPMaxPort {
		log.Fatalf("HTTPS port must be in range 0..%d", TCPMaxPort)
	}

	if Config.Verbose < 2 {
		return
	}

	log.Println("-- conf start ------------")
	log.Printf("MaxAge  = %d\n", Config.MaxAge)
	log.Printf("Verbose = %d\n", Config.Verbose)
	log.Printf("Http  -> { Address = %s, Port = %d }\n",
		Config.Http.Address, Config.Http.Port)
	log.Printf("Https -> { Address = %s, Port = %d }\n",
		Config.Https.Address, Config.Https.Port)
	log.Println("-- conf end --------------")
}

func ConfInit() {
	cf := flag.String("c", "", "JSON config file")
	ma := flag.Int("m", HTTPMaxAge, "content cache age in secs")
	ha := flag.String("a", "", "http address (default '' = all)")
	hp := flag.Int("p", 80, "http port")
	sa := flag.String("A", "", "https address (default '' = all)")
	sp := flag.Int("P", 443, "https port")
	vl := flag.Int("v", 0, "verbose level")
	flag.Parse()

	Config.MaxAge = *ma
	Config.Verbose = *vl
	Config.Http.Address = *ha
	Config.Http.Port = *hp
	Config.Https.Address = *sa
	Config.Https.Port = *sp
	defer ConfValidate()

	if *cf == "" {
		return
	}

	jf, err := os.Open(*cf)
	if err != nil {
		log.Println("unable to open conf file " + *cf)
		log.Fatal(err)
	}
	defer jf.Close()

	jsonParser := json.NewDecoder(jf)
	err = jsonParser.Decode(&Config)
	if err != nil {
		log.Println("error parsing JSON conf file " + *cf)
		log.Fatal(err)
	}
}

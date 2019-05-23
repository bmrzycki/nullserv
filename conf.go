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

var Config ConfFile

func ConfInit() {
	conf := flag.String("c", "", "JSON config file")
	ma := flag.Int("m", 31536000, "content cache age in secs")
	ha := flag.String("a", "", "http address (default '' = all)")
	hp := flag.Int("p", 80, "http port")
	hsa := flag.String("A", "", "https address (default '' = all)")
	hsp := flag.Int("P", 443, "https port")
	v := flag.Int("v", 0, "verbose level")
	flag.Parse()

	Config.MaxAge = *ma
	Config.Verbose = *v
	Config.Http.Address = *ha
	Config.Http.Port = *hp
	Config.Https.Address = *hsa
	Config.Https.Port = *hsp

	if *conf == "" {
		return
	}

	jf, err := os.Open(*conf)
	if err != nil {
		log.Println("unable to open conf file " + *conf)
		log.Fatal(err)
	}
	defer jf.Close()

	jsonParser := json.NewDecoder(jf)
	err = jsonParser.Decode(&Config)
	if err != nil {
		log.Println("error parsing JSON conf file " + *conf)
		log.Fatal(err)
	}

	if Config.Verbose > 1 {
		log.Println(Config)
	}
}

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/jacobwpeng/gotyrant"
)

type Config struct {
	Addr     string
	Key      string
	Filename string
}

func ce(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var config Config

func init() {
	flag.StringVar(&config.Filename, "file", "", "Data file name")
	flag.StringVar(&config.Addr, "addr", "", "TT address")
	flag.StringVar(&config.Key, "key", "", "Target key")
	flag.Parse()
}

func main() {
	client, err := gotyrant.NewClient(gotyrant.Config{
		Addr:    config.Addr,
		Timeout: time.Second * 3,
	})
	ce(err)
	data, err := client.Get([]byte(config.Key))
	ce(err)
	err = ioutil.WriteFile(config.Filename, data, 0644)
	ce(err)
	log.Printf("Retrieve value of %s at %s to %s done",
		config.Key, config.Addr, config.Filename)
}

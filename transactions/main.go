package main

import (
	"flag"
	"log"
	"os"
)

var (
	confPath   = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	fundsToken = flag.String("funds-token", os.Getenv("FUNDS_API_TOKEN"), "token to the funds api")
)

func main() {
	flag.Parse()

	if *confPath == "" {
		log.Fatal("Config path is empty")
	}

	if *fundsToken == "" {
		log.Fatal("No way to reach the funds service, no token provided")
	}
}

// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements a Promisance server.
package main

import (
	"log"
	"time"
)

func main() {
	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	tzName, _ := time.Now().Zone()
	log.Printf("app: server time zone is %s, logs are UTC\n", tzName)

	php, err := newInstance()
	if err != nil {
		log.Fatal(err)
	}

	if err := php.index_php(); err != nil {
		log.Fatalf("php: index: %v\n", err)
	}
}

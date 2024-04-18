// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements a Promisance server.
package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/mdhender/promisance/app/way"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	s := &server{host: "localhost", port: "8080"}
	s.tz, _ = time.Now().Zone()
	// developing, so use constant values
	s.setup.do = true
	s.setup.route = "/setup-375eb078-ae06-41f7-8258-b78372382c32"

	flag.BoolVar(&s.setup.do, "setup", false, "run setup")
	flag.StringVar(&s.data, "data", "", "path to store data files")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatal("error: unexpected values found on command line\n")
	}
	// verify data path
	if s.data = strings.TrimSpace(s.data); s.data == "" {
		log.Fatal("error: no data path specified\n")
	} else if path, err := filepath.Abs(s.data); err != nil {
		log.Fatalf("error: data: %v\n", err)
	} else if sb, err := os.Stat(path); err != nil {
		log.Fatalf("error: data: %s: no such directory\n", s.data)
	} else if !sb.IsDir() {
		log.Fatalf("error: data: %s: not a directory\n", s.data)
	} else {
		s.data = path
	}
	s.addr = net.JoinHostPort(s.host, s.port)

	log.Printf("app: server time zone is %s (logs are UTC)\n", s.tz)
	log.Printf("app: setup mode is %v\n", s.setup.do)
	log.Printf("app: data path  is %s\n", s.data)

	if s.setup.do {
		log.Printf("app: setup mode is enabled\n")
		// if we're running setup, we must verify that the directory is writeable
		log.Printf("app: verifying that data path is writable...\n")
		if file, err := os.CreateTemp(s.data, "example"); err != nil {
			log.Fatalf("error: data: %v\n", err)
		} else if err = os.Remove(file.Name()); err != nil { // cleanup
			log.Fatalf("error: data: %s\n", s.data)
		} else {
			log.Printf("app: data path is writable\n")
		}
	} else {
		php, err := newInstance()
		if err != nil {
			log.Fatal(err)
		}

		if err := php.index_php(); err != nil {
			log.Fatalf("php: index: %v\n", err)
		}
	}

	handler := s.routes()

	if s.setup.route == "" {
		log.Printf("app: serving on http://%s\n", s.addr)
	} else {
		log.Printf("app: serving setup on http://%s%s\n", s.addr, s.setup.route)
	}
	log.Fatalln(http.ListenAndServe(s.addr, handler))
}

type server struct {
	data  string // path to store data files
	addr  string
	host  string
	port  string
	setup struct {
		do    bool
		route string
	}
	tz string
}

func (s *server) routes() http.Handler {
	router := way.NewRouter()
	if s.setup.do {
		router.Handle("GET", "/", s.handleSetupIndex())
		// create a random router to serve the setup
		if s.setup.route == "" { // allow for a constant path for development
			s.setup.route = "/setup-" + uuid.New().String()
		}
		router.Handle("GET", s.setup.route, s.handleSetup())
		// we have to be very secure when running setup
		router.NotFound = s.handleNotAuthorized()
		return router
	}
	router.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello, World!")
	})
	return router
}

func (s *server) handleNotAuthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func (s *server) handleSetupIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<h1>?</h1><p>The server is down for maintenance. Please check back later.</p>`))
	}
}

func (s *server) handleSetup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

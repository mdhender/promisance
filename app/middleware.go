// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"log"
	"net/http"
	"strings"
)

// check_banned_ip
func (s *server) checkBannedIP() func(http.Handler) http.Handler {
	log.Printf("todo: implement checkedBannedIP middleware\n")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// define constants based on round start/end times
func (s *server) setRoundTimes() func(http.Handler) http.Handler {
	log.Printf("todo: implement setRoundTimes middleware\n")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//if CUR_TIME < s.world.round_time_begin { // pre-registration
			//	//define('ROUND_SIGNUP', TRUE);
			//	//define('ROUND_STARTED', FALSE);
			//	//define('ROUND_CLOSING', FALSE);
			//	//define('ROUND_FINISHED', FALSE);
			//	//define('TXT_TIMENOTICE', lang('ROUND_WILL_BEGIN', gmdate(lang('ROUND_WILL_BEGIN_FORMAT'), s.world.round_time_begin - CUR_TIME)));
			//} else if CUR_TIME < s.world.round_time_closing { // normal gameplay
			//	//define('ROUND_SIGNUP', TRUE);
			//	//define('ROUND_STARTED', TRUE);
			//	//define('ROUND_CLOSING', FALSE);
			//	//define('ROUND_FINISHED', FALSE);
			//} else if CUR_TIME < s.world.round_time_end { // final week (or so)
			//	//define('ROUND_SIGNUP', FALSE);
			//	//define('ROUND_STARTED', TRUE);
			//	//define('ROUND_CLOSING', TRUE);
			//	//define('ROUND_FINISHED', FALSE);
			//	//define('TXT_TIMENOTICE', lang('ROUND_WILL_END', gmdate(lang('ROUND_WILL_END_FORMAT'), s.world.round_time_end - CUR_TIME)));
			//} else { // end of round
			//	//define('ROUND_SIGNUP', FALSE);
			//	//define('ROUND_STARTED', FALSE);
			//	//define('ROUND_CLOSING', FALSE);
			//	//define('ROUND_FINISHED', TRUE);
			//	//define('TXT_TIMENOTICE', lang('ROUND_HAS_ENDED'));
			//}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *server) turnsCrontab() func(http.Handler) http.Handler {
	log.Printf("todo: implement turnsCrontab middleware\n")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func (s *server) validate_location(rules map[string]int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			page := strings.TrimPrefix(r.URL.Path, "/")
			rule, ok := rules[page]
			if !ok {
				// no rules to validate for this page
				log.Printf("middleware: %s %s: error.badpage\n", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
			}
			log.Printf("middleware: validate_location: page %s: rule %d: referer %q\n", page, rule, r.Referer())

			// pages that work from anywhere
			switch rule {
			case 0: // no referer needed
				// special case - go to "login" page when you still have an active session
				if page == "relogin" {
					log.Printf("middleware: %s %s: %d: redirect %s => main\n", r.Method, r.URL.Path, rule, page)
				}
			case 1: // need referer
				if r.Referer() == "" {
					log.Printf("middleware: %s %s: %d: error.noref\n", r.Method, r.URL.Path, rule)
				} else if page == "game" {
					log.Printf("middleware: %s %s: %d: redirect %s => game\n", r.Method, r.URL.Path, rule, page)
				}
			case 2: // need in-game referer
				if r.Referer() == "" {
					log.Printf("middleware: %s %s: %d: error.noref\n", r.Method, r.URL.Path, rule)
				} else if !strings.Contains(r.Referer(), URL_BASE) {
					log.Printf("middleware: %s %s: %d: error.badref\n", r.Method, r.URL.Path, rule)
				}
			default:
				log.Printf("middleware: %s %s: %d: error.badrule\n", r.Method, r.URL.Path, rule)
			}
			next.ServeHTTP(w, r)
		})
	}
}

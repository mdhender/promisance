// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"log"
	"net/http"
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

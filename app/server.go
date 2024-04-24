// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/promisance/app/authn"
	"github.com/mdhender/promisance/app/cerr"
	"github.com/mdhender/promisance/app/jot"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
	"log"
	"net/http"
	"strings"
)

type server struct {
	data            string // path to store data files
	templates       string // path to template files
	public          string // path to public files (sometimes called "assets")
	addr            string
	host            string
	port            string
	tz              string
	baseURL         string
	db              *orm.DB
	world           *model.World_t
	valid_locations map[string]int
	sessions        *jot.Factory_t
	authenticator   *authn.Authenticator
	language        *LanguageManager_t
}

func (s *server) check_banned_ip(ip string) bool {
	// todo: implement this check
	return false
}

func (s *server) checkAuth(r *http.Request, relogin bool) string {
	user := s.sessions.User(r)

	if !user.IsAuthenticated() {
		// user has no session cookie set
		if relogin {
			return s.language.Printf("ERROR_LOGIN_NO_SESSION") + "<br /><br />" + s.language.Printf("ERROR_PLEASE_RELOGIN_CHECK")
		}
		return s.language.Printf("ERROR_LOGIN_NO_SESSION")
	}

	log.Printf("%s %s: todo: checkAuth is assuming caller sets user1 and emp1 and language\n")

	return ""
}

// fetches a field from the posted form, trims whitespace
func (s *server) getFormVar(r *http.Request, key, defaultValue string) (string, bool) {
	values := r.URL.Query()[key]
	if len(values) == 0 {
		// try post
		if err := r.ParseForm(); err == nil {
			values = r.Form[key]
		}
	}
	if len(values) == 0 {
		return defaultValue, false
	}
	return strings.TrimSpace(values[0]), true
}

// log an event into the database
func (s *server) logmsg(kind PHPLoggingConstants, msg string) {
	log.Printf("todo: implement logmsg: %q\n", msg)
}

func (s *server) validate_location(r *http.Request, page string) (string, error) {
	rule, ok := s.valid_locations[page]
	if !ok {
		return "badpage", cerr.ErrBadPage
	}
	log.Printf("%s %s: validate_location: page %s: rule %d: referer %q\n", r.Method, r.URL, page, rule, r.Referer())

	// pages that work from anywhere
	switch rule {
	case 0: // no referer needed
		// special case - go to "login" page when you still have an active session
		if page == "relogin" {
			log.Printf("%s %s: validate_location: page %s => %q\n", r.Method, r.URL, page, "main")
			return "main", nil
		}
		log.Printf("%s %s: validate_location: page %s: rule %d => %q\n", r.Method, r.URL, page, rule, page)
		return page, nil
	case 1: // need referer
		if r.Referer() == "" {
			log.Printf("%s %s: validate_location: page %s: rule %d: missing referer\n", r.Method, r.URL, page, rule)
			return "noref", cerr.ErrMissingReferrer
		} else if page == "game" {
			log.Printf("%s %s: validate_location: page %s => %q\n", r.Method, r.URL, page, "main")
			return "main", nil
		}
		log.Printf("%s %s: validate_location: page %s: rule %d => %q\n", r.Method, r.URL, page, rule, page)
		return page, nil
	case 2: // need in-game referer
		if r.Referer() == "" {
			log.Printf("%s %s: validate_location: page %s: rule %d: missing referer\n", r.Method, r.URL, page, rule)
			return "noref", cerr.ErrMissingReferrer
		} else if !strings.Contains(r.Referer(), s.baseURL) {
			log.Printf("%s %s: validate_location: page %s: rule %d: bad referer %q\n", r.Method, r.URL, page, rule, r.Referer())
			return "badref", cerr.ErrBadReferrer
		}
		log.Printf("%s %s: validate_location: page %s: rule %d => %q\n", r.Method, r.URL, page, rule, page)
		return page, nil
	}
	log.Printf("%s %s: validate_location: page %s: rule %d: rule not found\n", r.Method, r.URL, page, rule)
	return "badpage", cerr.ErrBadPage
}

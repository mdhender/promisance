// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/promisance/app/jot"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginContent struct {
	GAME_TITLE       template.HTML
	LOGIN_VERSION    template.HTML
	LOGIN_DATE_RANGE template.HTML
	LOGIN_COUNTER    template.HTML
	NOTICES          template.HTML
	LABEL_USERNAME   template.HTML
	LABEL_PASSWORD   template.HTML
	LOGIN_SUBMIT     string
	SignupStatus     template.HTML
	LOGIN_TOPEMPIRES template.HTML
	CLAN_ENABLE      bool
	LOGIN_TOPCLANS   template.HTML
	LOGIN_TOPPLAYERS template.HTML
	LOGIN_HISTORY    template.HTML
	LOGIN_GUIDE      template.HTML
}

// loginGetHandler creates a login form.
// If the request has a valid session, redirect to "/signup" or "/game/{empireId}".
func (s *server) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL)
	sess := s.sessions.Session(r.Context())
	log.Printf("%s %s: session %p\n", r.Method, r.URL.Path, sess)
	log.Printf("%s %s: session %+v\n", r.Method, r.URL.Path, *sess)

	if sess.IsInvalid() {
		log.Printf("%s %s: session invalid: => /relogin\n", r.Method, r.URL.Path)
		http.Redirect(w, r, "/relogin", http.StatusSeeOther)
		return
	} else if sess.IsExpired() {
		log.Printf("%s %s: session expired: => /relogin\n", r.Method, r.URL.Path)
		http.Redirect(w, r, "/relogin", http.StatusSeeOther)
		return
	} else if sess.IsValid() {
		log.Printf("%s %s: session valid: empireId %d\n", r.Method, r.URL.Path, sess.empireId)
		if sess.empireId == 0 {
			log.Printf("%s %s: session valid: empireId %d: => /signup\n", r.Method, r.URL.Path, sess.empireId)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		log.Printf("%s %s: session valid: empireId %d: => /game/%d\n", r.Method, r.URL.Path, sess.empireId, sess.empireId)
		http.Redirect(w, r, fmt.Sprintf("/game/%d", sess.empireId), http.StatusSeeOther)
		return
	}
	log.Printf("%s %s: session missing\n", r.Method, r.URL.Path)

	// do the form
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Login</title><link rel="stylesheet" href="https://unpkg.com/missing.css@1.1.1"></head><body>`))
	_, _ = w.Write([]byte(`<h1>Login</h1>`))
	_, _ = w.Write([]byte(`<main>`))
	_, _ = w.Write([]byte(`<p>Insert login page here</p>`))
	_, _ = w.Write([]byte(`<form method="post" action="/login" class="box rows">`))
	_, _ = w.Write([]byte(`<p>`))
	_, _ = w.Write([]byte(`    <label for="login_username">{{.LABEL_USERNAME}}</label>`))
	_, _ = w.Write([]byte(`    <input type="text" name="login_username" size="18" id="login_username" value="basque"/>`))
	_, _ = w.Write([]byte(`</p>`))
	_, _ = w.Write([]byte(`<p>`))
	_, _ = w.Write([]byte(`    <label for="login_username">{{.LABEL_PASSWORD}}</label>`))
	_, _ = w.Write([]byte(`    <input type="password" name="login_password" size="18" id="login_password" value="bisque"/>`))
	_, _ = w.Write([]byte(`</p>`))
	_, _ = w.Write([]byte(`<input type="hidden" name="action" value="login" />`))
	_, _ = w.Write([]byte(`<input type="submit" value="{{.LOGIN_SUBMIT}}"/>`))
	_, _ = w.Write([]byte(`</form>`))
	_, _ = w.Write([]byte(`<ol>`))
	_, _ = w.Write([]byte(`<li><a href="/">Home</a></li>`))
	_, _ = w.Write([]byte(`<li><a href="/relogin">Relogin</a></li>`))
	_, _ = w.Write([]byte(`</ol>`))
	_, _ = w.Write([]byte(`</main>`))
	_, _ = w.Write([]byte(`</body>`))
}

func (s *server) loginGetHandlerOld(w http.ResponseWriter, r *http.Request) {
	started := time.Now()

	// if the session cookie is set, redirect to the main page.
	// if the session is invalid, that page will delete it and redirect them back to here
	user := s.jots.User(r)
	log.Printf("%s %s: lgh: user %+v\n", r.Method, r.URL, user)
	log.Printf("%s %s: lgh: user: authenticated %v\n", r.Method, r.URL, user.IsAuthenticated())
	if user.IsAuthenticated() {
		http.Redirect(w, r, "/relogin", http.StatusSeeOther)
		return
	}

	// explicitly destroy the session, which will clear the token cookie
	s.jots.Destroy(w)

	// our response variables
	content := LoginContent{
		GAME_TITLE:       GAME_TITLE,
		LOGIN_VERSION:    s.language.PrintfHTML("LOGIN_VERSION", GAME_VERSION),
		LOGIN_DATE_RANGE: s.language.PrintfHTML("LOGIN_DATE_RANGE", s.world.RoundTimeBegin, s.world.RoundTimeEnd),
		NOTICES:          s.noticesFromQueryParameter(r, 1),
		LABEL_USERNAME:   s.language.PrintfHTML("LABEL_USERNAME"),
		LABEL_PASSWORD:   s.language.PrintfHTML("LABEL_PASSWORD"),
		LOGIN_SUBMIT:     s.language.Printf("LOGIN_SUBMIT"),
		LOGIN_TOPEMPIRES: template.HTML(fmt.Sprintf(`<a href="/topempires"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPEMPIRES"))),
		CLAN_ENABLE:      CLAN_ENABLE,
		LOGIN_TOPCLANS:   template.HTML(fmt.Sprintf(`<a href="/topclans"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPCLANS"))),
		LOGIN_TOPPLAYERS: template.HTML(fmt.Sprintf(`<a href="/topplayers"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPPLAYERS"))),
		LOGIN_HISTORY:    template.HTML(fmt.Sprintf(`<a href="/history"><b>%s</b></a><br />`, s.language.Printf("LOGIN_HISTORY"))),
		LOGIN_GUIDE:      template.HTML(fmt.Sprintf(`<a href="/guide"><b>%s</b></a><br />`, s.language.Printf("LOGIN_GUIDE"))),
	}
	if num, err := s.db.EmpireActiveCount(); err != nil {
		log.Printf("%s %s: empireActiveCount: %v\n", r.Method, r.URL, err)
		content.LOGIN_COUNTER = "<b>***</b>"
	} else {
		content.LOGIN_COUNTER = template.HTML(fmt.Sprintf("<b>%03d</b>", num))
	}
	if COUNTER_TEMPLATE != "" {
		log.Printf("%s %s: counter template is not implemented\n", r.Method, r.URL)
		//counter, err := getimagesize(filepath.Join(PROM_BASEDIR, "images", COUNTER_TEMPLATE))
		//if err != nil {
		//	log.Printf("%s %s: lgh: error getting image size: %v\n", r.Method, r.URL, err)
		//} else {
		//	countData = fmt.Sprintf(`<img src="?location=count" alt="%s" style="width:%dpx;height:%dpx" />`, countData, counter[0]/10*len(countData), counter[1])
		//}
	}
	if ROUND_SIGNUP && !(SIGNUP_CLOSED_USER && SIGNUP_CLOSED_EMPIRE) {
		content.SignupStatus = template.HTML(fmt.Sprintf(`<a href="/signup"><b>%s</b></a><br />`, s.language.Printf("LOGIN_SIGNUP")))
	} else {
		content.SignupStatus = template.HTML(fmt.Sprintf(`<b>%s</b><br />`, s.language.Printf("LOGIN_SIGNUP_CLOSED")))
	}

	layout := CompactLayoutPayload{
		Header:  s.getCompactHeader("login"),
		Content: content,
		Footer:  s.getCompactFooter(started),
	}

	// render the login page template
	log.Printf("%s %s: lgh: rendering from template\n", r.Method, r.URL)
	s.render(w, r, layout, "html_compact.gohtml", "login.gohtml")
}

func (s *server) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL)

	// Get the form values
	username := r.FormValue("login_username")
	log.Printf("%s %s: login_username: %q\n", r.Method, r.URL, username)
	password := r.FormValue("login_password")
	log.Printf("%s %s: login_password: %q", r.Method, r.URL, password)

	// Validate the form inputs
	var notices []string
	if username == "" {
		notices = append(notices, s.language.Printf("INPUT_NEED_USERNAME"))
	} else if strings.TrimSpace(username) != username {
		notices = append(notices, "Username must not start or end with spaces.")
	}
	if password == "" {
		notices = append(notices, s.language.Printf("INPUT_NEED_PASSWORD"))
	} else if strings.TrimSpace(password) != password {
		notices = append(notices, "Password must not start or end with spaces.")
	}
	if len(notices) != 0 {
		log.Printf("%s %s: notices %v\n", r.Method, r.URL, notices)
		args, ok := s.noticesToQueryParameters(notices)
		if !ok {
			log.Printf("%s %s: notices: => /login\n", r.Method, r.URL)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Printf("%s %s: notices: => /logi?%sn\n", r.Method, r.URL, args)
		http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
		return
	}

	log.Printf("%s %s: authenticated\n", r.Method, r.URL)

	sess, err := s.sessions.Create(1, 0)
	if err != nil {
		log.Printf("%s %s: sessions: create %v\n", r.Method, r.URL, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sess.CreateCookie(w)

	log.Printf("%s %s: authenticated: => /\n", r.Method, r.URL)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) loginPostHandlerOld(w http.ResponseWriter, r *http.Request) {
	// Get the form values
	username := r.FormValue("login_username")
	log.Printf("%s %s: login_username: %q\n", r.Method, r.URL, username)
	password := r.FormValue("login_password")
	log.Printf("%s %s: login_password: %q", r.Method, r.URL, password)

	// Validate the form inputs
	var notices []string
	if username == "" {
		notices = append(notices, s.language.Printf("INPUT_NEED_USERNAME"))
	} else if strings.TrimSpace(username) != username {
		notices = append(notices, "Username must not start or end with spaces.")
	}
	if password == "" {
		notices = append(notices, s.language.Printf("INPUT_NEED_PASSWORD"))
	} else if strings.TrimSpace(password) != password {
		notices = append(notices, "Password must not start or end with spaces.")
	}
	if len(notices) != 0 {
		log.Printf("%s %s: notices %v\n", r.Method, r.URL, notices)
		args, ok := s.noticesToQueryParameters(notices)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
		return
	}

	// Authenticate the user
	user, err := s.authenticator.Authenticate(username, password)
	if err != nil {
		log.Printf("%s %s: lph: authentication failed\n", r.Method, r.URL)
		args, ok := s.noticesToQueryParameters([]string{"Invalid credentials"})
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
		return
	}
	log.Printf("%s %s: lph: authentication succeded\n", r.Method, r.URL)
	roles := s.authenticator.UserRoles(user)
	roles["authenticated"] = true
	jUser := jot.User_t{UserId: user.Id, EmpireId: 1, Roles: roles}
	log.Printf("%s %s: user %v\n", r.Method, r.URL, user)

	if user.UserName == "" {
		// this should be impossible if authentication succeeded
		notices = []string{s.language.Printf("LOGIN_USER_NOT_FOUND")}
		s.logmsg(E_USER_NOTICE, "failed (load) - "+username)
		args, ok := s.noticesToQueryParameters([]string{""})
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
		return
	}

	if user.Flags.Closed {
		notices = []string{s.language.Printf("LOGIN_USER_CLOSED")}
		s.logmsg(E_USER_NOTICE, "failed (closed) - "+username)
		args, ok := s.noticesToQueryParameters([]string{""})
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
		return
	}

	// Retrieve the associated empires
	empList, err := s.db.UserActiveEmpires(user.Id)
	if err != nil {
		log.Printf("%s %s: lph: error retrieving empires: %v\n", r.Method, r.URL, err)
		log.Printf("%s %s: %s\n", r.Method, r.URL, s.language.Printf("ERROR_TITLE", "Failed to check for registered empires"))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if the user has any empires. If they do, fetch just the first one.
	if len(empList) == 0 {
		log.Printf("%s %s: empList is empty\n", r.Method, r.URL)
		// if they've signed up before but don't have an empire, bounce them over to the signup page
		if !ROUND_SIGNUP {
			log.Printf("%s %s: empList is empty: ROUND_SIGNUP is false\n", r.Method, r.URL)
			if args, ok := s.noticesToQueryParameters([]string{s.language.Printf("LOGIN_NO_EMPIRE")}); ok {
				http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if SIGNUP_CLOSED_EMPIRE {
			log.Printf("%s %s: empList is empty: SIGNUP_CLOSED_EMPIRE is true\n", r.Method, r.URL)
			if args, ok := s.noticesToQueryParameters([]string{s.language.Printf("LOGIN_NO_EMPIRE_CLOSED")}); ok {
				http.Redirect(w, r, "/login?"+args, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Printf("%s %s: redirecting before saving session cookies\n", r.Method, r.URL)
		http.Redirect(w, r, "/signup&registered="+username, http.StatusSeeOther)
		return
	}

	// use the first empire owned by this user when creating the session
	jUser.EmpireId = empList[0].Id
	cookie, err := s.jots.NewTokenCookie(7*24*time.Hour, jUser)
	if err != nil {
		log.Printf("%s %s: lph: sessions token failed: %v\n", r.Method, r.URL, err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, cookie)
	log.Printf("%s %s: session cookie %v\n", r.Method, r.URL, cookie)

	// only set them online if the round has actually started
	if ROUND_STARTED {
		empList[0].Flags.Online = true
		if err := s.db.EmpireUpdateFlags(empList[0]); err != nil {
			log.Printf("%s %s: empireUpdateFlags failed: %v\n", r.Method, r.URL, err)
		}
	}

	// Update the user's last IP and last date
	user.LastIP, user.LastDate = r.RemoteAddr, time.Now().UTC()
	if err := s.db.UserAccessUpdate(user); err != nil {
		log.Printf("%s %s: userAccessUpdate failed: %v\n", r.Method, r.URL, err)
	}

	// Redirect to the game location
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

// logoutGetHandler deletes any session cookies and redirects back to the home page.
func (s *server) logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	sess := s.sessions.Session(r.Context())
	log.Printf("%s %s: session %p\n", r.Method, r.URL.Path, sess)
	log.Printf("%s %s: session %+v\n", r.Method, r.URL.Path, *sess)
	if sess != nil && sess.id != "" {
		s.sessions.DestroySession(sess.id)
	}
	s.sessions.DestroyCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// logoutPostHandler deletes any session cookies and redirects back to the home page.
func (s *server) logoutPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	sess := s.sessions.Session(r.Context())
	log.Printf("%s %s: session %p\n", r.Method, r.URL.Path, sess)
	log.Printf("%s %s: session %+v\n", r.Method, r.URL.Path, *sess)
	if sess != nil && sess.id != "" {
		s.sessions.DestroySession(sess.id)
	}
	s.sessions.DestroyCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// reloginGetHandler deletes any session cookies and redirects back to the home page.
func (s *server) reloginGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	sess := s.sessions.Session(r.Context())
	log.Printf("%s %s: session %p\n", r.Method, r.URL.Path, sess)
	log.Printf("%s %s: session %+v\n", r.Method, r.URL.Path, *sess)
	if sess != nil && sess.id != "" {
		s.sessions.DestroySession(sess.id)
	}
	s.sessions.DestroyCookies(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/promisance/app/jot"
	"html/template"
	"log"
	"net/http"
	"time"
)

type LoginContent struct {
	GAME_TITLE       template.HTML
	LOGIN_VERSION    template.HTML
	LOGIN_DATE_RANGE template.HTML
	LOGIN_COUNTER    template.HTML
	Notices          []string // populate with 0 or 1, no more
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

func (s *server) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	started := time.Now()

	// redirect to the main page if the user is authenticated
	user := s.sessions.User(r)
	log.Printf("%s %s: lgh: user %+v\n", r.Method, r.URL, user)
	if user.IsAuthenticated() {
		log.Printf("%s %s: lgh: user is authenticated\n", r.Method, r.URL)
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}
	log.Printf("%s %s: lgh: user is not authenticated\n", r.Method, r.URL)
	// explicitly clear the token cookie
	log.Printf("%s %s: lgh: deleted cookies\n", r.Method, r.URL)
	s.sessions.DeleteCookie(w)

	// our response variables
	content := LoginContent{
		GAME_TITLE:       GAME_TITLE,
		LOGIN_VERSION:    s.language.PrintfHTML("LOGIN_VERSION", GAME_VERSION),
		LOGIN_DATE_RANGE: s.language.PrintfHTML("LOGIN_DATE_RANGE", s.world.RoundTimeBegin, s.world.RoundTimeEnd),
		LABEL_USERNAME:   s.language.PrintfHTML("LABEL_USERNAME"),
		LABEL_PASSWORD:   s.language.PrintfHTML("LABEL_PASSWORD"),
		LOGIN_SUBMIT:     s.language.Printf("LOGIN_SUBMIT"),
		LOGIN_TOPEMPIRES: template.HTML(fmt.Sprintf(`<a href="/index.php?location=topempires"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPEMPIRES"))),
		CLAN_ENABLE:      CLAN_ENABLE,
		LOGIN_TOPCLANS:   template.HTML(fmt.Sprintf(`<a href="/index.php?location=topclans"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPCLANS"))),
		LOGIN_TOPPLAYERS: template.HTML(fmt.Sprintf(`<a href="/index.php?location=topplayers"><b>%s</b></a><br />`, s.language.Printf("LOGIN_TOPPLAYERS"))),
		LOGIN_HISTORY:    template.HTML(fmt.Sprintf(`<a href="/index.php?location=history"><b>%s</b></a><br />`, s.language.Printf("LOGIN_HISTORY"))),
		LOGIN_GUIDE:      template.HTML(fmt.Sprintf(`<a href="/index.php?location=guide"><b>%s</b></a><br />`, s.language.Printf("LOGIN_GUIDE"))),
	}
	num := 3 // $db->queryCell('SELECT COUNT(*) FROM '. EMPIRE_TABLE .' WHERE u_id != 0');
	content.LOGIN_COUNTER = template.HTML(fmt.Sprintf("<b>%03d</b>", num))
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
		content.SignupStatus = template.HTML(fmt.Sprintf(`<a href="/index.php?location=signup"><b>%s</b></a><br />`, s.language.Printf("LOGIN_SIGNUP")))
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
	// Get the form values
	username := r.FormValue("login_username")
	log.Printf("%s %s: login_username: %q\n", r.Method, r.URL, username)
	password := r.FormValue("login_password")
	log.Printf("%s %s: login_password: %q", r.Method, r.URL, password)

	// Validate the form inputs
	// ...

	// Authenticate the user
	if !s.authenticator.Authenticate(username, password) {
		log.Printf("%s %s: lph: authentication failed\n", r.Method, r.URL)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	log.Printf("%s %s: lph: authentication succeded\n", r.Method, r.URL)
	user := jot.User_t{UserId: 1, EmpireId: 1, Roles: map[string]bool{"authenticated": true}}
	log.Printf("%s %s: user %v\n", r.Method, r.URL, user)

	// create a new token and save it as a cookie
	cookie, err := s.sessions.NewTokenCookie(7*24*time.Hour, user)
	if err != nil {
		log.Printf("%s %s: lph: sessions token failed: %v\n", r.Method, r.URL, err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, cookie)

	//	// Retrieve the associated empires
	//	empires, err := h.getEmpires(user.ID)
	//	if err != nil {
	//		// Handle error retrieving empires
	//		// ...
	//		return
	//	}
	//
	//	// Check if the user has any empires
	//	if len(empires) == 0 {
	//		// Handle the case where the user has no empires
	//		// ...
	//		return
	//	}
	//
	//	// Load the first empire
	//	empire := empires[0]
	//
	//	// Initialize the session
	//	err = h.session.Start(w, r)
	//	if err != nil {
	//		// Handle session initialization error
	//		// ...
	//		return
	//	}
	//
	//	// Set the user and empire in the session
	//	h.session.Set("user", user)
	//	h.session.Set("empire", empire)
	//
	//	// Update the user's last IP and last date
	//	user.LastIP = r.RemoteAddr
	//	user.LastDate = time.Now()
	//
	//	// Save the user and empire
	//	err = h.db.SaveUser(user)
	//	if err != nil {
	//		// Handle error saving user
	//		// ...
	//		return
	//	}
	//
	//	err = h.db.SaveEmpire(empire)
	//	if err != nil {
	//		// Handle error saving empire
	//		// ...
	//		return
	//	}

	// Redirect to the game location
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

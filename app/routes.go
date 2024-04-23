// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mdhender/promisance/app/jot"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

func (s *server) routes(valid_locations map[string]int) http.Handler {
	if s.sessions == nil {
		panic("assert(sessions != nil)")
	}

	router := chi.NewRouter()

	// public routes, no authentication required, okay to cache
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/site-map", http.StatusTemporaryRedirect)
		})
		r.Get("/site-map", s.sitemapHandler)
	})

	// login/logout pages, no authentication required, do not cache
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/login", s.loginGetHandler)
		r.Post("/login", s.loginPostHandler)
		r.Get("/logout", s.logoutHandler)
	})

	// , s.checkBannedIP(), s.validate_location(valid_locations), s.turnsCrontab(), s.setRoundTimes())

	// admin pages, authentication required, do not cache
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger, middleware.NoCache, s.sessions.Authenticator())
		r.Get("/admin/clans", s.adminClansHandler)
		r.Get("/admin/empedit", s.adminEmpeditHandler)
		r.Get("/admin/empires", s.adminEmpiresHandler)
		r.Get("/admin/history", s.adminHistoryHandler)
		r.Get("/admin/log", s.adminLogHandler)
		r.Get("/admin/market", s.adminMarketHandler)
		r.Get("/admin/messages", s.adminMessagesHandler)
		r.Get("/admin/permissions", s.adminPermissionsHandler)
		r.Get("/admin/round", s.adminRoundHandler)
		r.Get("/admin/users", s.adminUsersHandler)
	})

	// player pages, authentication required, do not cache
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger, middleware.NoCache, s.sessions.Authenticator())
		r.Get("/aid", s.aidHandler)
		r.Get("/bank", s.bankHandler)
		r.Get("/banner", s.bannerHandler)
		r.Get("/build", s.buildHandler)
		r.Get("/cash", s.cashHandler)
		r.Get("/clan", s.clanHandler)
		r.Get("/clanforum", s.clanforumHandler)
		r.Get("/clanstats", s.clanstatsHandler)
		r.Get("/contacts", s.contactsHandler)
		r.Get("/count", s.countHandler)
		r.Get("/credits", s.creditsHandler)
		r.Get("/delete", s.deleteHandler)
		r.Get("/demolish", s.demolishHandler)
		r.Get("/farm", s.farmHandler)
		r.Get("/game", s.gameHandler)
		r.Get("/graveyard", s.graveyardHandler)
		r.Get("/guide", s.guideHandler)
		r.Get("/history", s.historyHandler)
		r.Get("/land", s.landHandler)
		r.Get("/lottery", s.lotteryHandler)
		r.Get("/magic", s.magicHandler)
		r.Get("/main", s.mainHandler)
		r.Get("/manage/clan", s.manageClansHandler)
		r.Get("/manage/empire", s.manageEmpireHandler)
		r.Get("/manage/user", s.manageUserHandler)
		r.Get("/messages", s.messagesHandler)
		r.Get("/military", s.militaryHandler)
		r.Get("/news", s.newsHandler)
		r.Get("/pguide", s.pguideHandler)
		r.Get("/playerstats", s.playerstatsHandler)
		r.Get("/pubmarketbuy", s.pubmarketbuyHandler)
		r.Get("/pubmarketsell", s.pubmarketsellHandler)
		r.Get("/pvtmarketbuy", s.pvtmarketbuyHandler)
		r.Get("/pvtmarketsell", s.pvtmarketsellHandler)
		r.Get("/relogin", s.reloginHandler)
		r.Get("/revalidate", s.revalidateHandler)
		r.Get("/scores", s.scoresHandler)
		r.Get("/search", s.searchHandler)
		r.Get("/signup", s.signupHandler)
		r.Get("/status", s.statusHandler)
		r.Get("/topclans", s.topclansHandler)
		r.Get("/topempires", s.topempiresHandler)
		r.Get("/topplayers", s.topplayersHandler)
		r.Get("/validate", s.validateHandler)
	})

	router.NotFound(s.assetsHandler(s.public))

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
		// php handlers can panic, so catch and deal with them
		defer func() {
			if catch := recover(); catch != nil {
				log.Printf("php: panic: %v\n\n%s\n", catch, debug.Stack())
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`<h1>Setup Panic!</h1><pre>`))
				_, _ = w.Write(debug.Stack())
				_, _ = w.Write([]byte(`</pre>`))
				return
			}
		}()

		log.Printf("app: setup mode is enabled\n")
		// if we're running setup, we must verify that the directory is writeable
		log.Printf("app: verifying that data path is writable...\n")
		if file, err := os.CreateTemp(s.data, "example"); err != nil {
			log.Printf("setup: data: not writable\n")
			log.Printf("setup: data: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if err = os.Remove(file.Name()); err != nil { // cleanup
			log.Printf("setup: data: not writable\n")
			log.Printf("setup: data: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Printf("setup: data %s: writable\n", s.data)

		php, err := newInstance(w, r)
		if err != nil {
			log.Fatal(err)
		}
		if err := php.install_setup_php(); err != nil {
			log.Fatalf("php: index: %v\n", err)
		}

		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *server) adminClansHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminEmpeditHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminEmpiresHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminHistoryHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminLogHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminMarketHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminMessagesHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminRoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) aidHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) bankHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) bannerHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) buildHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) cashHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) clanHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) clanforumHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) clanstatsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) contactsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) countHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) creditsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) demolishHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) farmHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) gameHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) graveyardHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) guideHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) historyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) landHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

type CompactLayoutPayload struct {
	Header  *CompactHeaderPayload
	Content any
	Footer  *CompactFooterPayload
}
type CompactHeaderPayload struct {
	Page       string // internal name of the page?
	Title      string
	LANG_CODE  string
	LANG_DIR   string
	GetStyles  string
	AddStyles  []string
	AddScripts []string
}
type CompactFooterPayload struct {
	HTML_FOOTER       template.HTML
	HTML_LINK_CREDITS string
	HTML_LINK_LOGIN   string
	DEBUG_FOOTER      bool
	HTML_DEBUG_FOOTER template.HTML
}

func (s *server) getCompactHeader(page string) *CompactHeaderPayload {
	return &CompactHeaderPayload{
		Page:      page,
		Title:     s.language.Printf("HTML_TITLE", "login"),
		LANG_CODE: s.language.Printf("LANG_CODE"),
		LANG_DIR:  s.language.Printf("LANG_DIR"),
		GetStyles: "qmt.css",
	}
}
func (s *server) getCompactFooter(started time.Time) *CompactFooterPayload {
	dur, memUsage, peakMemUsage, queryCount := time.Now().Sub(started), 1, 2, 3
	return &CompactFooterPayload{
		HTML_FOOTER:       s.language.PrintfHTML("HTML_FOOTER", GAME_VERSION),
		HTML_LINK_CREDITS: s.language.Printf("HTML_LINK_CREDITS"),
		HTML_LINK_LOGIN:   s.language.Printf("HTML_LINK_LOGIN"),
		DEBUG_FOOTER:      true,
		HTML_DEBUG_FOOTER: s.language.PrintfHTML("HTML_DEBUG_FOOTER", dur, memUsage, peakMemUsage, queryCount),
	}
}

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
func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.sessions.DeleteCookie(w)
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) lotteryHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) magicHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) manageClansHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) manageEmpireHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) manageUserHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) messagesHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) militaryHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) newsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) pguideHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) playerstatsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) pubmarketbuyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) pubmarketsellHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) pvtmarketbuyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) pvtmarketsellHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) reloginHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) revalidateHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) scoresHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) searchHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) signupHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) sitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html><head><title>Promisance Site Map</title></head><body>
<h1>Promisance Site Map</h1>
<table>
<tr><td><code>/admin/clans</code></td><td><a href="/admin/clans">Admin Clans Handler</a></td></tr>
<tr><td><code>/admin/empedit</code></td><td><a href="/admin/empedit">Admin Empire Edit Handler</a></td></tr>
<tr><td><code>/admin/empires</code></td><td><a href="/admin/empires">Admin Empires Handler</a></td></tr>
<tr><td><code>/admin/history</code></td><td><a href="/admin/history">Admin History Handler</a></td></tr>
<tr><td><code>/admin/log</code></td><td><a href="/admin/log">Admin Log Handler</a></td></tr>
<tr><td><code>/admin/market</code></td><td><a href="/admin/market">Admin Market Handler</a></td></tr>
<tr><td><code>/admin/messages</code></td><td><a href="/admin/messages">Admin Messages Handler</a></td></tr>
<tr><td><code>/admin/permissions</code></td><td><a href="/admin/permissions">Admin Permissions Handler</a></td></tr>
<tr><td><code>/admin/round</code></td><td><a href="/admin/round">Admin Round Handler</a></td></tr>
<tr><td><code>/admin/users</code></td><td><a href="/admin/users">Admin Users Handler</a></td></tr>
<tr><td><code>/aid</code></td><td><a href="/aid">Aid Handler</a></td></tr>
<tr><td><code>/bank</code></td><td><a href="/bank">Bank Handler</a></td></tr>
<tr><td><code>/banner</code></td><td><a href="/banner">Banner Handler</a></td></tr>
<tr><td><code>/build</code></td><td><a href="/build">Build Handler</a></td></tr>
<tr><td><code>/cash</code></td><td><a href="/cash">Cash Handler</a></td></tr>
<tr><td><code>/clan</code></td><td><a href="/clan">Clan Handler</a></td></tr>
<tr><td><code>/clanforum</code></td><td><a href="/clanforum">Clan Forum Handler</a></td></tr>
<tr><td><code>/clanstats</code></td><td><a href="/clanstats">Clan Stats Handler</a></td></tr>
<tr><td><code>/contacts</code></td><td><a href="/contacts">Contacts Handler</a></td></tr>
<tr><td><code>/count</code></td><td><a href="/count">Count Handler</a></td></tr>
<tr><td><code>/credits</code></td><td><a href="/credits">Credits Handler</a></td></tr>
<tr><td><code>/delete</code></td><td><a href="/delete">Delete Handler</a></td></tr>
<tr><td><code>/demolish</code></td><td><a href="/demolish">Demolish Handler</a></td></tr>
<tr><td><code>/farm</code></td><td><a href="/farm">Farm Handler</a></td></tr>
<tr><td><code>/game</code></td><td><a href="/game">Game Handler</a></td></tr>
<tr><td><code>/graveyard</code></td><td><a href="/graveyard">Graveyard Handler</a></td></tr>
<tr><td><code>/guide</code></td><td><a href="/guide">Guide Handler</a></td></tr>
<tr><td><code>/history</code></td><td><a href="/history">History Handler</a></td></tr>
<tr><td><code>/land</code></td><td><a href="/land">Land Handler</a></td></tr>
<tr><td><code>/login</code></td><td><a href="/login">Login Handler</a></td></tr>
<tr><td><code>/logout</code></td><td><a href="/logout">Logout Handler</a></td></tr>
<tr><td><code>/lottery</code></td><td><a href="/lottery">Lottery Handler</a></td></tr>
<tr><td><code>/magic</code></td><td><a href="/magic">Magic Handler</a></td></tr>
<tr><td><code>/main</code></td><td><a href="/main">Main Handler</a></td></tr>
<tr><td><code>/manage/clan</code></td><td><a href="/manage/clan">Manage Clans Handler</a></td></tr>
<tr><td><code>/manage/empire</code></td><td><a href="/manage/empire">Manage Empire Handler</a></td></tr>
<tr><td><code>/manage/user</code></td><td><a href="/manage/user">Manage User Handler</a></td></tr>
<tr><td><code>/messages</code></td><td><a href="/messages">Messages Handler</a></td></tr>
<tr><td><code>/military</code></td><td><a href="/military">Military Handler</a></td></tr>
<tr><td><code>/news</code></td><td><a href="/news">News Handler</a></td></tr>
<tr><td><code>/pguide</code></td><td><a href="/pguide">Player's Guide Handler</a></td></tr>
<tr><td><code>/playerstats</code></td><td><a href="/playerstats">Player Stats Handler</a></td></tr>
<tr><td><code>/pubmarketbuy</code></td><td><a href="/pubmarketbuy">Public Market Buy Handler</a></td></tr>
<tr><td><code>/pubmarketsell</code></td><td><a href="/pubmarketsell">Public Market Sell Handler</a></td></tr>
<tr><td><code>/pvtmarketbuy</code></td><td><a href="/pvtmarketbuy">Private Market Buy Handler</a></td></tr>
<tr><td><code>/pvtmarketsell</code></td><td><a href="/pvtmarketsell">Private Market Sell Handler</a></td></tr>
<tr><td><code>/relogin</code></td><td><a href="/relogin">Relogin Handler</a></td></tr>
<tr><td><code>/revalidate</code></td><td><a href="/revalidate">Revalidate Handler</a></td></tr>
<tr><td><code>/scores</code></td><td><a href="/scores">Scores Handler</a></td></tr>
<tr><td><code>/search</code></td><td><a href="/search">Search Handler</a></td></tr>
<tr><td><code>/signup</code></td><td><a href="/signup">Signup Handler</a></td></tr>
<tr><td><code>/status</code></td><td><a href="/status">Status Handler</a></td></tr>
<tr><td><code>/topclans</code></td><td><a href="/topclans">Top Clans Handler</a></td></tr>
<tr><td><code>/topempires</code></td><td><a href="/topempires">Top Empires Handler</a></td></tr>
<tr><td><code>/topplayers</code></td><td><a href="/topplayers">Top Players Handler</a></td></tr>
<tr><td><code>/validate</code></td><td><a href="/validate">Validate Handler</a></td></tr>
</table>
</body></html>`))
}
func (s *server) statusHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) topclansHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) topempiresHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) topplayersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) validateHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *server) render(w http.ResponseWriter, r *http.Request, payload any, templates ...string) {
	var files []string
	for _, t := range templates {
		files = append(files, filepath.Join(s.templates, t))
	}

	var err error
	t, err := template.New("layout").Funcs(template.FuncMap{
		"yield": func() (string, error) {
			return "", fmt.Errorf("yield called unexpectedly.")
		},
	}).ParseFiles(files...)
	if err != nil {
		log.Printf("%s %s: template: parse: %v", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, "layout", payload); err != nil {
		log.Printf("%s %s: template: execute: %v", r.Method, r.URL, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (s *server) assetsHandler(assetsPath string) http.HandlerFunc {
	cacheControl := fmt.Sprintf("public, max-age=%d, immutable", 28*24*60*60)
	log.Printf("server: assets: cache-control %q\n", cacheControl)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Join the assets path with the url path.
		// Join calls path.Clean on the result for us automatically.
		path := filepath.Join(assetsPath, r.URL.Path)

		// check whether a file exists or is a directory at the given path
		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// todo: set cache control header to serve file for a month or so
		// static files in this case need to be cache busted
		// (usually by appending a hash to the filename)
		w.Header().Set("Cache-Control", cacheControl)

		// otherwise, use http.FileServer to serve the asset.

		// we're creating a new file server on every request.
		// that server deals with the request and then goes away.
		http.FileServer(http.Dir(assetsPath)).ServeHTTP(w, r)
	}
}

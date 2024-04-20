// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/promisance/app/way"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func (s *server) routes() http.Handler {
	router := way.NewRouter()
	router.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/site-map", http.StatusTemporaryRedirect)
	})
	router.HandleFunc("GET", "/admin/clans", s.adminClansHandler)
	router.HandleFunc("GET", "/admin/empedit", s.adminEmpeditHandler)
	router.HandleFunc("GET", "/admin/empires", s.adminEmpiresHandler)
	router.HandleFunc("GET", "/admin/history", s.adminHistoryHandler)
	router.HandleFunc("GET", "/admin/log", s.adminLogHandler)
	router.HandleFunc("GET", "/admin/market", s.adminMarketHandler)
	router.HandleFunc("GET", "/admin/messages", s.adminMessagesHandler)
	router.HandleFunc("GET", "/admin/permissions", s.adminPermissionsHandler)
	router.HandleFunc("GET", "/admin/round", s.adminRoundHandler)
	router.HandleFunc("GET", "/admin/users", s.adminUsersHandler)
	router.HandleFunc("GET", "/aid", s.aidHandler)
	router.HandleFunc("GET", "/bank", s.bankHandler)
	router.HandleFunc("GET", "/banner", s.bannerHandler)
	router.HandleFunc("GET", "/build", s.buildHandler)
	router.HandleFunc("GET", "/cash", s.cashHandler)
	router.HandleFunc("GET", "/clan", s.clanHandler)
	router.HandleFunc("GET", "/clanforum", s.clanforumHandler)
	router.HandleFunc("GET", "/clanstats", s.clanstatsHandler)
	router.HandleFunc("GET", "/contacts", s.contactsHandler)
	router.HandleFunc("GET", "/count", s.countHandler)
	router.HandleFunc("GET", "/credits", s.creditsHandler)
	router.HandleFunc("GET", "/delete", s.deleteHandler)
	router.HandleFunc("GET", "/demolish", s.demolishHandler)
	router.HandleFunc("GET", "/farm", s.farmHandler)
	router.HandleFunc("GET", "/game", s.gameHandler)
	router.HandleFunc("GET", "/graveyard", s.graveyardHandler)
	router.HandleFunc("GET", "/guide", s.guideHandler)
	router.HandleFunc("GET", "/history", s.historyHandler)
	router.HandleFunc("GET", "/land", s.landHandler)
	router.HandleFunc("GET", "/login", s.loginHandler)
	router.HandleFunc("GET", "/logout", s.logoutHandler)
	router.HandleFunc("GET", "/lottery", s.lotteryHandler)
	router.HandleFunc("GET", "/magic", s.magicHandler)
	router.HandleFunc("GET", "/main", s.mainHandler)
	router.HandleFunc("GET", "/manage/clan", s.manageClansHandler)
	router.HandleFunc("GET", "/manage/empire", s.manageEmpireHandler)
	router.HandleFunc("GET", "/manage/user", s.manageUserHandler)
	router.HandleFunc("GET", "/messages", s.messagesHandler)
	router.HandleFunc("GET", "/military", s.militaryHandler)
	router.HandleFunc("GET", "/news", s.newsHandler)
	router.HandleFunc("GET", "/pguide", s.pguideHandler)
	router.HandleFunc("GET", "/playerstats", s.playerstatsHandler)
	router.HandleFunc("GET", "/pubmarketbuy", s.pubmarketbuyHandler)
	router.HandleFunc("GET", "/pubmarketsell", s.pubmarketsellHandler)
	router.HandleFunc("GET", "/pvtmarketbuy", s.pvtmarketbuyHandler)
	router.HandleFunc("GET", "/pvtmarketsell", s.pvtmarketsellHandler)
	router.HandleFunc("GET", "/relogin", s.reloginHandler)
	router.HandleFunc("GET", "/revalidate", s.revalidateHandler)
	router.HandleFunc("GET", "/scores", s.scoresHandler)
	router.HandleFunc("GET", "/search", s.searchHandler)
	router.HandleFunc("GET", "/signup", s.signupHandler)
	router.HandleFunc("GET", "/site-map", s.sitemapHandler)
	router.HandleFunc("GET", "/status", s.statusHandler)
	router.HandleFunc("GET", "/topclans", s.topclansHandler)
	router.HandleFunc("GET", "/topempires", s.topempiresHandler)
	router.HandleFunc("GET", "/topplayers", s.topplayersHandler)
	router.HandleFunc("GET", "/validate", s.validateHandler)

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
func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
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

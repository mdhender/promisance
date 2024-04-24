// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mdhender/promisance/app/cerr"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/way"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

func (s *server) routes() http.Handler {
	r := way.NewRouter()
	r.Handle("GET", "/index.php", s.indexPhpHandler())
	r.NotFound = s.assetsHandler(s.public)
	if r != nil {
		return r
	}

	if s.sessions == nil {
		panic("assert(sessions != nil)")
	}

	router := chi.NewRouter()

	// public routes, no authentication required, okay to cache
	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//	http.Redirect(w, r, "/site-map", http.StatusTemporaryRedirect)
		//})
		r.Get("/", s.indexPhpHandler())
		r.Get("/index.php", s.indexPhpHandler())
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

type sessionVariables_t struct {
	User      *model.User_t
	Empire    *model.Empire_t
	World     *model.World_t
	Round     model.RoundData_t
	Started   time.Time
	Title     string
	NeedPriv  model.UserFlag_t
	Locks     map[string]int
	Action    string
	Language  string
	LogMsg    string
	Page      string
	TriedPage string
}

type sessionVariablesContextKey_t string

func (s *server) indexPhpHandler() http.HandlerFunc {
	// temporarily save some routing information. we don't use it, but may.
	// Valid in-game pages - can be specified for 'location' parameter to load corresponding PHP file
	// Values denote any special requirements for loading the page
	if s.valid_locations == nil {
		s.valid_locations = map[string]int{
			// 0 - does not require referer or session
			"count":       0, // 0 - does not require referer or session
			"credits":     0, // 0 - does not require referer or session
			"history":     0, // 0 - does not require referer or session
			"login":       0, // 0 - does not require referer or session
			"pguide":      0, // 0 - does not require referer or session
			"playerstats": 0, // 0 - does not require referer or session
			"signup":      0, // 0 - does not require referer or session
			"topclans":    0, // 0 - does not require referer or session
			"topempires":  0, // 0 - does not require referer or session
			"topplayers":  0, // 0 - does not require referer or session
			"relogin":     0, // redirect from login page load; redirects don't set referer, and this could be a bookmark

			// 1 - requires referer from any site
			"game": 1, // redirect from login page submission; redirects don't set referer

			// 2 - requires referer from in-game, also requires active session
			"banner":     2, // 2 - requires referer from in-game, also requires active session
			"guide":      2, // 2 - requires referer from in-game, also requires active session
			"messages":   2, // 2 - requires referer from in-game, also requires active session
			"revalidate": 2, // 2 - requires referer from in-game, also requires active session
			"validate":   2, // 2 - requires referer from in-game, also requires active session
			"main":       2, // both "relogin" and "game" redirect to here

			// Information
			"clanstats": 2, // 2 - requires referer from in-game, also requires active session
			"contacts":  2, // 2 - requires referer from in-game, also requires active session
			"graveyard": 2, // 2 - requires referer from in-game, also requires active session
			"news":      2, // 2 - requires referer from in-game, also requires active session
			"scores":    2, // 2 - requires referer from in-game, also requires active session
			"search":    2, // 2 - requires referer from in-game, also requires active session
			"status":    2, // 2 - requires referer from in-game, also requires active session

			// Use Turns
			"build":    2, // 2 - requires referer from in-game, also requires active session
			"cash":     2, // 2 - requires referer from in-game, also requires active session
			"demolish": 2, // 2 - requires referer from in-game, also requires active session
			"farm":     2, // 2 - requires referer from in-game, also requires active session
			"land":     2, // 2 - requires referer from in-game, also requires active session

			// Finances
			"bank":          2, // 2 - requires referer from in-game, also requires active session
			"lottery":       2, // 2 - requires referer from in-game, also requires active session
			"pubmarketbuy":  2, // 2 - requires referer from in-game, also requires active session
			"pubmarketsell": 2, // 2 - requires referer from in-game, also requires active session
			"pvtmarketbuy":  2, // 2 - requires referer from in-game, also requires active session
			"pvtmarketsell": 2, // 2 - requires referer from in-game, also requires active session

			// Foreign Affairs
			"aid":       2, // 2 - requires referer from in-game, also requires active session
			"clan":      2, // 2 - requires referer from in-game, also requires active session
			"clanforum": 2, // 2 - requires referer from in-game, also requires active session
			"magic":     2, // 2 - requires referer from in-game, also requires active session
			"military":  2, // 2 - requires referer from in-game, also requires active session

			// Management
			"delete":        2, // 2 - requires referer from in-game, also requires active session
			"manage/clan":   2, // 2 - requires referer from in-game, also requires active session
			"manage/empire": 2, // 2 - requires referer from in-game, also requires active session
			"manage/user":   2, // 2 - requires referer from in-game, also requires active session

			// Administration
			"admin/clans":       2, // 2 - requires referer from in-game, also requires active session
			"admin/empedit":     2, // 2 - requires referer from in-game, also requires active session
			"admin/empires":     2, // 2 - requires referer from in-game, also requires active session
			"admin/history":     2, // 2 - requires referer from in-game, also requires active session
			"admin/log":         2, // 2 - requires referer from in-game, also requires active session
			"admin/market":      2, // 2 - requires referer from in-game, also requires active session
			"admin/messages":    2, // 2 - requires referer from in-game, also requires active session
			"admin/permissions": 2, // 2 - requires referer from in-game, also requires active session
			"admin/round":       2, // 2 - requires referer from in-game, also requires active session
			"admin/users":       2, // 2 - requires referer from in-game, also requires active session

			// Logout
			"logout": 2, // 2 - requires referer from in-game, also requires active session
		}
		log.Printf("todo: implement valid_locations referer logic (%d pages)\n", len(s.valid_locations))
	}

	showTurnsCrontabError := true

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered\n", r.Method, r.URL)

		sv := &sessionVariables_t{
			Started: time.Now(),
		}

		if s.db == nil {
			log.Printf("%s %s: ERROR_TITLE %s\n", r.Method, r.URL, s.language.Printf("ERROR_TITLE"))
			log.Printf("%s %s: ERROR_DATABASE_OFFLINE %s\n", r.Method, r.URL, s.language.Printf("ERROR_DATABASE_OFFLINE"))
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		// load world variables
		if s.world == nil {
			log.Printf("%s %s: world not initialized\n", r.Method, r.URL)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		// If we're configured for cron-less turn updates, check them now
		// turns := fetchSomeTurnData()
		if !TURNS_CRONTAB {
			if showTurnsCrontabError {
				log.Printf("%s %s: turns crontab is not implemented\n", r.Method, r.URL)
				showTurnsCrontabError = false
			}
			// turns.doUpdate()
		}

		// define constants based on round start/end times
		if sv.Started.Before(s.world.RoundTimeBegin) { // pre-registration
			sv.Round.Signup = true
			sv.Round.Started = false
			sv.Round.Closing = false
			sv.Round.Finished = false
			sv.Round.TimeNotice = s.language.Printf("ROUND_WILL_BEGIN", "ROUND_WILL_BEGIN_FORMAT", s.world.RoundTimeBegin.Sub(sv.Started))
		} else if sv.Started.Before(s.world.RoundTimeClosing) { // normal gameplay
			sv.Round.Signup = true
			sv.Round.Started = true
			sv.Round.Closing = false
			sv.Round.Finished = false
		} else if sv.Started.Before(s.world.RoundTimeEnd) { // final week (or so)
			sv.Round.Signup = false
			sv.Round.Started = true
			sv.Round.Closing = true
			sv.Round.Finished = false
			sv.Round.TimeNotice = s.language.Printf("ROUND_WILL_END", "ROUND_WILL_BEGIN_FORMAT", s.world.RoundTimeEnd.Sub(sv.Started))
		} else { // end of round
			sv.Round.Signup = false
			sv.Round.Started = false
			sv.Round.Closing = false
			sv.Round.Finished = true
			sv.Round.TimeNotice = s.language.Printf("ROUND_HAS_ENDED")
		}
		// todo: inject round data into request context

		if s.check_banned_ip("ip.address") {
			log.Printf("%s %s: implement banned ip address logic\n", r.Method, r.URL)
			// $ban_message = lang('YOU_ARE_BANNED', gmdate(lang('COMMON_TIME_FORMAT'), $ban['p_createtime']), ($ban['p_reason']) ? $ban['p_reason'] : lang('BANNED_NO_REASON'), ($ban['p_expire'] == 0) ? lang('BANNED_PERMANENT') : lang('BANNED_EXPIRES', gmdate(lang('COMMON_TIME_FORMAT'), $ban['p_expire'])), MAIL_ADMIN);
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Special variables parsed by page_header()
		// Page title ("Promisance - whatever")

		// Set to combination of UFLAG_MOD/UFLAG_ADMIN to indicate USER privileges required to load page
		// var needpriv model.UserFlag_t

		// Add entries to this array to request entities to be loaded and locked.
		sv.Locks = map[string]int{"emp1": 0, "emp2": 0, "user1": 0, "user2": 0, "clan1": 0, "clan2": 0, "world": 0}
		// Set to an entity number, or use -1 (for non-empires) to determine the ID automatically from the loaded empire
		// If loading an entity fails, its value will be reset to 0
		// Setting 'emp1' has no effect - it exists for logging purposes and is automatically set to the current empire ID
		// Setting 'user1' or 'world' to anything other than -1 has no effect - they only allow auto-detection.
		// Additional locks (for special purposes) can be requested directly from $db

		sv.TriedPage, _ = s.getFormVar(r, "location", "login")
		var errchk error
		sv.Page, errchk = s.validate_location(r, sv.TriedPage)
		sv.Action, _ = s.getFormVar(r, "action", "")

		// if they tried entering a really, really long action, truncate it and log a warning
		if len(sv.Action) > 64 {
			s.logmsg(E_USER_NOTICE, "action overflowed: "+sv.Action)
			sv.Action = sv.Action[:64]
		}

		if errchk != nil {
			message := "<table><tr><th>" + s.language.Printf("SECURITY_TITLE") + "</th></tr><tr><td>" + s.language.Printf("SECURITY_DESC") + "<br />"
			error_args := []string{"triedpage", "action"}
			sv.LogMsg = sv.TriedPage
			if errors.Is(errchk, cerr.ErrBadPage) {
				message += s.language.Printf("SECURITY_BADPAGE", url.QueryEscape(sv.TriedPage)) + "<br />"
			} else if errors.Is(errchk, cerr.ErrBadReferrer) {
				referer := r.Referer()
				message += s.language.Printf("SECURITY_BADREF", url.QueryEscape(sv.TriedPage)+url.QueryEscape(referer)+s.baseURL) + "<br />"
				error_args = append(error_args, "referer")
			} else if errors.Is(errchk, cerr.ErrMissingReferrer) {
				message += s.language.Printf("SECURITY_NOREF", url.QueryEscape(sv.TriedPage)) + "<br />"
			} else {
				message += s.language.Printf("SECURITY_UNKNOWN", url.QueryEscape(sv.TriedPage), errchk) + "<br />"
			}
			message += s.language.Printf("SECURITY_INSTRUCT", s.baseURL, MAIL_ADMIN) + "</td></tr></table>"
			s.logmsg(E_USER_ERROR, "varlist($error_args, get_defined_vars())")

			if errors.Is(errchk, cerr.ErrBadPage) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			} else {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}

		// check if destination page requires an active session
		if s.valid_locations[sv.Page] == 2 {
			// prevent logout page from suggesting to log back in again
			auth := s.checkAuth(r, sv.Page != "logout")
			if auth != "" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}

		// logic to set user1, language, and emp1 pulled from checkAuth()
		sUser := s.sessions.User(r)
		var err error
		if sUser.IsAuthenticated() {
			sv.User, err = s.db.UserFetch(sUser.UserId)
			if err != nil {
				log.Printf("%s %s: userFetch %v\n", r.Method, r.URL, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			sv.Language = sv.User.Lang
			sv.Empire, err = s.db.EmpireFetch(sUser.EmpireId)
			if err != nil {
				log.Printf("%s %s: empireFetch %v\n", r.Method, r.URL, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		// finally, finally, finally...
		ctx := r.Context()
		ctx = context.WithValue(ctx, sessionVariablesContextKey_t("session_variables"), sv)

		// dispatch to the page handler
		switch sv.Page {
		case "login":
			log.Printf("%s %s: routing to page %s\n", r.Method, r.URL, sv.Page)
			switch r.Method {
			case "GET":
				s.loginGetHandler(w, r.WithContext(ctx))
			case "POST":
				http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
			default:
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
		default:
			log.Printf("%s %s: page %s: not implemented\n", r.Method, r.URL, sv.Page)
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		}

		//location := r.URL.Query().Get("location")
		//if location == "" {
		//	log.Printf("%s %s: no location parameter\n", r.Method, r.URL)
		//	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		//	return
		//}
		//log.Printf("%s %s: location %q\n", r.Method, r.URL, location)
		//scheme := "http"
		//if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		//	scheme = "https"
		//}
		//locationURL := fmt.Sprintf("%s://%s/%s", scheme, r.Host, location)
		//log.Printf("%s %s: redirecting to %s\n", r.Method, r.URL, locationURL)
		//http.Redirect(w, r, locationURL, http.StatusSeeOther)
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

func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.sessions.Destroy(w)
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
func (s *server) noticesToQueryParameters(notices []string) (string, bool) {
	var parms string
	for _, notice := range notices {
		if parms != "" {
			parms += "&"
		}
		parms += "notice=" + base64.URLEncoding.EncodeToString([]byte(notice))
	}
	if len(parms) == 0 {
		return "", false
	}
	return parms, true
}
func (s *server) noticesFromQueryParameter(r *http.Request, style int) template.HTML {
	log.Printf("%s %s: nfqp: style %d\n", r.Method, r.URL.Path, style)
	args := r.URL.Query()["notice"]
	log.Printf("%s %s: nfqp: args  %v\n", r.Method, r.URL.Path, args)
	if len(args) == 0 {
		return ""
	}
	var notices []string
	for n, arg := range args {
		log.Printf("%s %s: nfqp: arg   %d %q\n", r.Method, r.URL.Path, n, arg)
		if len(arg) != 0 {
			if msg, err := base64.URLEncoding.DecodeString(arg); err == nil && len(msg) != 0 {
				log.Printf("%s %s: nfqp: msg   %d %q\n", r.Method, r.URL.Path, n, string(msg))
				notices = append(notices, string(msg))
			}
		}
	}
	return s.notices(style, notices)
}

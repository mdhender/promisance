// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"github.com/mdhender/promisance/app/orm"
	"log"
	"net/http"
	"strings"
	"time"
)

type sessionStore_t struct {
	store *orm.DB
	// defaults are applied when the store returns an empty session_t
	defaults struct {
		ttl  time.Duration
		lang string
	}
}

type sessionContext_t string

func NewSessionStore(store *orm.DB, ttl time.Duration, lang string) *sessionStore_t {
	return &sessionStore_t{
		store: store,
		defaults: struct {
			ttl  time.Duration
			lang string
		}{
			ttl:  ttl,
			lang: lang,
		},
	}
}

// Authenticator is middleware that retrieves a session_t from the request, and adds it to the request's context.
// If token is not found or is invalid, an empty session_t is added instead.
func (s *sessionStore_t) Authenticator(next http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("sessions: authenticator: entered\n")
		started := time.Now()

		var sess *session_t
		id := s.sessionIdFromRequest(r)
		if id == "" {
			log.Printf("%s %s: sessions: authenticator: no session found\n", r.Method, r.URL)
			sess = &session_t{}
		} else {
			log.Printf("%s %s: sessions: authenticator: session %s: found\n", r.Method, r.URL, id)
			// extract the session information from the session store
			sess = s.sessionFromStore(id)
			log.Printf("%s %s: sessions: authenticator: store %+v\n", r.Method, r.URL.Path, *sess)
		}

		// add value to the session by setting default values and timing
		sess.started = started
		if sess.lang == "" {
			sess.lang = s.defaults.lang
		}

		log.Printf("%s %s: sessions: authenticator: session %+v\n", r.Method, r.URL.Path, *sess)

		// add the session_t to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, sessionContext_t("session_t"), sess)

		// this is middleware, so move on to the next handler in the chain
		log.Printf("sessions: authenticator: calling next\n")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Create creates a new session in the store.
// Returns an error if unable to do so.
// Otherwise, returns a session_t with the new session data.
func (s *sessionStore_t) Create(userId, empireId int) (*session_t, error) {
	id, err := s.store.SessionCreate(userId, empireId, s.defaults.ttl)
	if err != nil {
		return nil, err
	}
	return &session_t{
		id:       id,
		userId:   userId,
		empireId: empireId,
		lang:     s.defaults.lang,
		started:  time.Now(),
	}, nil
}

func (s *sessionStore_t) DestroyCookies(w http.ResponseWriter) {
	log.Printf("sessions: destroying cookies\n")
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "session_t",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
}

func (s *sessionStore_t) DestroySession(id string) {
	log.Printf("sessions: destroying session %q\n", id)
	if id == "" {
		return
	}
	_ = s.store.SessionsPurgeId(id)
}

// Session retrieves the session_t from the context.
// Returns an empty session_t if no session is found.
func (s *sessionStore_t) Session(ctx context.Context) *session_t {
	if sess, ok := ctx.Value(sessionContext_t("session_t")).(*session_t); ok {
		return sess
	}
	return &session_t{lang: s.defaults.lang}
}

// sessionFromStore retrieves the session information from the store.
// Returns an empty session_t if no session is found.
func (s *sessionStore_t) sessionFromStore(id string) *session_t {
	data, err := s.store.SessionFetch(id)
	if err != nil {
		log.Printf("sessions: store: %s: fetch %v\n", id, err)
		return &session_t{invalid: true}
	}
	return &session_t{
		id:       id,
		userId:   data.UserID,
		empireId: data.EmpireID,
	}
}

// sessionFromRequest extracts a session ID from a request.
// It looks at the bearer token first, then at the cookie.
// Returns an empty string if no session ID is found.
func (s *sessionStore_t) sessionIdFromRequest(r *http.Request) string {
	if id := s.sessionIdFromBearerToken(r); id != "" {
		return id
	}
	return s.sessionIdFromCookie(r)
}

// sessionIdFromBearerToken extracts a session ID from the request's bearer token.
// Returns an empty string if no session ID is found.
func (s *sessionStore_t) sessionIdFromBearerToken(r *http.Request) string {
	log.Printf("sessions: bearer: entered\n")
	headerAuthText := r.Header.Get("Authorization")
	if headerAuthText == "" {
		return ""
	}
	log.Printf("sessions: bearer: found authorization header\n")
	authTokens := strings.SplitN(headerAuthText, " ", 2)
	if len(authTokens) != 2 {
		return ""
	}
	log.Printf("sessions: bearer: found authorization token\n")
	kind, id := authTokens[0], strings.TrimSpace(authTokens[1])
	if kind != "Bearer" {
		return ""
	}
	log.Printf("sessions: bearer: found %q\n", id)
	return id
}

// sessionFromCookie extracts a session_t from cookies in the request.
// Returns an empty string if no session ID is found.
func (s *sessionStore_t) sessionIdFromCookie(r *http.Request) string {
	log.Printf("sessions: cookie: entered\n")
	c, err := r.Cookie("session_t")
	if err != nil {
		log.Printf("sessions: cookie: %+v\n", err)
		return ""
	}
	log.Printf("sessions: cookie: found %q\n", c.Value)
	return c.Value
}

type session_t struct {
	id       string
	userId   int
	empireId int
	lang     string
	expired  bool
	invalid  bool
	started  time.Time
}

func (s *session_t) IsExpired() bool {
	return s != nil && s.expired
}

func (s *session_t) IsInvalid() bool {
	return s != nil && s.invalid
}

func (s *session_t) IsMissing() bool {
	return s == nil || s.id == ""
}

// IsValid returns true if the session exists, and it isn't invalid or expired.
func (s *session_t) IsValid() bool {
	return !s.IsMissing() && !(s.invalid || s.expired)
}

func (s *session_t) CreateCookie(w http.ResponseWriter) {
	log.Printf("sessions: creating cookies\n")
	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "session_t",
		Value:    s.id,
		Expires:  time.Now().Add(7 * 24 * time.Hour).UTC(),
		HttpOnly: true,
		Secure:   true,
	})
}

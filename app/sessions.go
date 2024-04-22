// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"time"
)

// Sessions is middleware that extracts sessions from requests and adds the session data
// to the context for use by handlers. If the request doesn't contain a session, an
// "unauthenticated" user with no roles is added.
func (s *server) Sessions(sm *SessionManager_t) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("sessions.mw: entered\n")
			// an uninitialized session user is treated as an unauthenticated user
			var session *Session_t

			// use the cookie value to fetch the associated session and update the session user
			if cookie, err := r.Cookie("sid"); err == nil {
				session = sm.GetSession(cookie.Value)
			}

			// add the session user to the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, sessionContextKey_t("session"), session)

			// this is middleware, so move on to the next handler in the chain
			log.Printf("sessions.mw: calling next\n")
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

type SessionManager_t struct {
	sync.Mutex
	lastSweep time.Time
	// sessions is a map of session_id to the session
	sessions map[string]*Session_t
}

// Session_t is a user's current session.
// A nil Session_t pointer is treated the same as an unauthenticated user.
type Session_t struct {
	Id       string // unique id for the session
	UserId   int    // user id
	EmpireId int    // current empire id
	Role     string
	Expires  time.Time
}

// NewSession creates a new session or updates an existing session.
// A user may only have one session.
// Warning: it has the side effect of culling expired sessions.
// Warning: may create race conditions since it updates session data.
func (sm *SessionManager_t) NewSession(user, empire int) *Session_t {
	sm.Lock()
	defer sm.Unlock()

	var session *Session_t

	// does the user currently have a session?
	for _, sess := range sm.sessions {
		if sess.UserId == user {
			session = sess
			break
		}
		if sess.isExpired() {
			delete(sm.sessions, sess.Id)
		}
	}

	// create a new session only if needed
	if session == nil {
		session = &Session_t{
			Id:     uuid.New().String(),
			UserId: user,
			Role:   "user",
		}
		sm.sessions[session.Id] = session
	}

	session.EmpireId = empire

	// default session time-to-live is two weeks
	session.Expires = time.Now().Add(2 * 7 * 24 * time.Hour)

	return session
}

func (sm *SessionManager_t) DeleteSession(id string) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.sessions, id)
}

// GetSession uses the id to retrieve a session.
// If the id isn't valid, a nil session pointer is returned.
// If the id is valid but the session is expired, it is deleted from
// the manager. (Nasty little side effect.) In that case, a nil session
// pointer is also returned.
func (sm *SessionManager_t) GetSession(id string) *Session_t {
	sm.Lock()
	defer sm.Unlock()

	session, ok := sm.sessions[id]
	if !ok { // no session found
		return nil
	} else if session.isExpired() {
		delete(sm.sessions, id)
		return nil
	}

	return session
}

// User returns the session user from in the supplied context.
// If there is none, an unauthenticated user is returned instead.
func (sm *SessionManager_t) User(ctx context.Context) SessionUser_t {
	session, ok := ctx.Value(sessionContextKey_t("session")).(*Session_t)
	if !ok || session == nil {
		return SessionUser_t{}
	}
	return SessionUser_t{
		UserId:   session.UserId,
		EmpireId: session.EmpireId,
	}
}

// CreateCookie is a method of the Session struct. This function sets an HTTP cookie
// with specified attributes using the standard Go library.
//
// Parameters:
//
//	w http.ResponseWriter: an HTTP response writer instance where we can set the cookie headers.
//
// Return:
//
//	This function doesn't return a value.
func (s *Session_t) CreateCookie(w http.ResponseWriter) {
	if s.isExpired() {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    s.Id,
		Expires:  s.Expires.UTC(),
		Path:     "/",  // Path for which the cookie is valid
		HttpOnly: true, // make sure HttpOnly is true to prevent javascript access
		Secure:   true, // make sure Secure is true if over HTTPS
	})
}

// DeleteCookie is a method of the Session struct. This function deletes an HTTP cookie.
//
// Parameters:
//
//	w http.ResponseWriter: an HTTP response writer instance where we can set the cookie headers.
//
// Return:
//
//	This function doesn't return a value.
func (s *Session_t) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		MaxAge:   -1,   // Set MaxAge to -1 to delete the cookie
		Path:     "/",  // Path for which the cookie is valid
		HttpOnly: true, // Ensure not accessible via Javascript
		Secure:   true, // Ensure sent over HTTPS
	})
}

func (s *Session_t) isAuthenticated() bool {
	return s != nil
}

func (s *Session_t) isAuthorized() bool {
	return s != nil && s.UserId != 0
}

func (s *Session_t) isExpired() bool {
	return s == nil || !time.Now().Before(s.Expires)
}

// sessionContextKey_t is the context key type for storing the Session in the context.
type sessionContextKey_t string

type SessionUser_t struct {
	UserId   int
	EmpireId int
}

func (su SessionUser_t) isAuthenticated() bool {
	return su.UserId != 0
}

func (su SessionUser_t) isAuthorized() bool {
	return su.UserId != 0
}

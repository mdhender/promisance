// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package jot implements a JOT with a payload that's tailored to this application.
package jot

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/mdhender/semver"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Version is the version of this package.
func Version() string {
	return semver.Version{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}.String()
}

// NewFactory returns an initialized factory.
// Name and Path are for cookies. Leave blank for default values.
// TTL is for all tokens.
// Signer is the initial signer.
func NewFactory(name, path string, ttl time.Duration, signer Signer_i) (*Factory_t, error) {
	if signer == nil {
		return nil, ErrInvalidSigner
	} else if signer.Expired() {
		return nil, ErrSignerExpired
	}
	f := Factory_t{
		ttl:     ttl,
		signers: map[string]Signer_i{signer.Id(): signer},
	}
	if name == "" {
		f.cookie.name = "promisance_jot"
	} else {
		f.cookie.name = name
	}
	if path == "" {
		f.cookie.path = "/"
	} else {
		f.cookie.path = path
	}
	return &f, nil
}

// Authenticator is middleware that retrieves a JOT from the request,
// extracts the user data, and adds it to the request's context.
// If token is not found or is invalid, the "unauthenticated" user is added instead.
func (f *Factory_t) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("jot: authenicator: entered\n")
			user, ok := f.payloadFromRequest(r)
			if !ok {
				log.Printf("jot: authenicator: no token found\n")
			}

			// add the session user to the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, userContextKey_t("user_t"), user)

			// this is middleware, so move on to the next handler in the chain
			log.Printf("jot: authenticator: calling next\n")
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func (f *Factory_t) User(r *http.Request) User_t {
	user, ok := r.Context().Value(userContextKey_t("user_t")).(User_t)
	if !ok {
		return unauthenticatedUser
	}
	return user

}

func (f *Factory_t) Destroy(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Path:     f.cookie.path,
		Name:     f.cookie.name,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Path:     f.cookie.path,
		Name:     "fh-auth",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Path:     f.cookie.path,
		Name:     "session_id",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
}

// userContextKey_t is the context key type for storing the jot.User_t in the context.
type userContextKey_t string

// JOT implements my version of the JSON Web JOT.
// It's like a JWT that's been customized for this application.
type JOT struct {
	Header    Header_t
	Claims    Claims_t
	Signature []byte
	isSigned  bool
}

// IsNotExpired returns true if the token has not expired.
func (j *JOT) IsNotExpired() bool {
	return j.Claims.IsNotExpired(time.Now().UTC())
}

// IsSigned returns true only if the signature has been verified.
func (j *JOT) IsSigned() bool {
	panic("!implemented")
}

// IsValid returns true only if the token is signed and not expired.
func (j *JOT) IsValid() bool {
	return !j.IsNotExpired() && j.IsSigned()
}

// String implements the Stringer interface.
func (j *JOT) String() string {
	panic("!implemented")
}

// Header_t is the header from a JOT.
type Header_t struct {
	Algorithm   string `json:"alg"`            // message authentication code algorithm, required
	ContentType string `json:"cty,omitempty"`  // not implemented
	Critical    string `json:"crit,omitempty"` // not implemented
	KeyID       string `json:"kid"`            // identifier used to sign, required
	TokenType   string `json:"typ"`            // should always be JOT, required
}

// Encode marshals the Header to JSON, then encodes the result as Base64.
func (h Header_t) Encode() ([]byte, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return encode_bytes(data), nil
}

// Claims_t is sometimes called the payload.
type Claims_t struct {
	// The recipients that the JOT is intended for.
	// Each principal intended to process the JOT must identify itself with a value in the audience claim.
	// If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present,
	// then the JOT must be rejected.
	// Not implemented.
	// Audience []string `json:"aud,omitempty"`

	// The expiration time on and after which the JOT must not be accepted for processing.
	ExpiresAt NumericDate_t `json:"exp"`

	// The time at which the JOT was issued.
	// Not implemented.
	IssuedAt NumericDate_t `json:"iat,omitempty"`

	// The principal that issued the JOT.
	// Not implemented.
	// Issuer string `json:"iss,omitempty"`

	// Case-sensitive unique identifier of the token even among different issuers.
	// Not implemented.
	// JWTID string `json:"jti,omitempty"`

	// The time on which the JOT will start to be accepted for processing.
	// Not implemented.
	// NotBefore NumericDate `json:"nbf,omitempty"`

	// The subject of the JOT.
	// Not implemented
	// Subject string `json:"sub"`

	Payload User_t `json:"payload"`
}

// Encode marshals the Claims to JSON, then encodes the result as Base64.
func (c Claims_t) Encode() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return encode_bytes(data), nil
}

// IsExpired returns true if `expiresAt` is not after `now`.
func (c Claims_t) IsExpired(now time.Time) bool {
	return !time.Time(c.ExpiresAt).After(now)
}

// IsNotExpired returns true if `expiresAt` is after `now`.
func (c Claims_t) IsNotExpired(now time.Time) bool {
	return time.Time(c.ExpiresAt).After(now)
}

// NumericDate_t is a numeric date representing seconds past 1970-01-01 00:00:00Z.
type NumericDate_t time.Time

// MarshalJSON implements the json.Marshaler interface.
func (d NumericDate_t) MarshalJSON() ([]byte, error) {
	n := time.Time(d).Unix()
	return json.Marshal(n)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *NumericDate_t) UnmarshalJSON(data []byte) error {
	var n int64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*d = NumericDate_t(time.Unix(n, 0))
	return nil
}

// String implements the Stringer interface.
func (d NumericDate_t) String() string {
	return time.Time(d).Format(time.RFC3339)
}

// User_t is the data for this user's session.
type User_t struct {
	UserId   int     `json:"user_id"`
	EmpireId int     `json:"empire_id,omitempty"`
	Roles    Roles_t `json:"roles,omitempty"`
}

func (u User_t) IsAuthenticated() bool {
	return u.Roles["authenticated"]
}

// Roles_t is a map of all roles assigned to a Subject.
type Roles_t map[string]bool

// MarshalJSON implements the json.Marshaler interface.
func (r Roles_t) MarshalJSON() ([]byte, error) {
	var roles []string
	for k := range r {
		roles = append(roles, k)
	}
	sort.Strings(roles)
	return json.Marshal(roles)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Roles_t) UnmarshalJSON(data []byte) error {
	var roles []string
	if err := json.Unmarshal(data, &roles); err != nil {
		return err
	}
	rr := make(map[string]bool)
	for _, role := range roles {
		rr[role] = true
	}
	*r = rr
	return nil
}

// decode_bytes is a helper for base-64 decoding.
func decode_bytes(src []byte) ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(src)))
	_, err := base64.RawURLEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// decode_str is a helper for base-64 decoding.
func decode_str(src string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(src)
}

// encode_bytes is a helper for base-64 encoding
func encode_bytes(src []byte) []byte {
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(dst, src)
	return dst
}

// encode_str is a helper for base-64 encoding
func encode_str(src string) []byte {
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(dst, []byte(src))
	return dst
}

// Errors used by the package.
const (
	ErrBadFactory       = constError("bad factory")
	ErrBadSigner        = constError("bad signer")
	ErrClaimsExpired    = constError("claims expired")
	ErrInvalidAlgorithm = constError("invalid algorithm")
	ErrInvalidHeader    = constError("invalid header")
	ErrInvalidSignature = constError("invalid signature")
	ErrInvalidSigner    = constError("invalid signer")
	ErrInvalidToken     = constError("invalid token")
	ErrMissingClaims    = constError("missing claims")
	ErrMissingSigner    = constError("missing signer")
	ErrNotFound         = constError("not found")
	ErrSignerExpired    = constError("signer expired")
	ErrUnauthorized     = constError("unauthorized")
	ErrUnknownType      = constError("unknown type")
)

// declarations to support constant errors
type constError string

func (ce constError) Error() string {
	return string(ce)
}

// tokenFromRequest extracts a token from a request.
// It looks at the bearer token first, then at the cookie.
// Returns an empty string if no token is found.
func (f *Factory_t) tokenFromRequest(r *http.Request) string {
	token := f.tokenFromBearerToken(r)
	if token == "" {
		token = f.tokenFromCookie(r)
	}
	return token
}

// tokenFromBearerToken extracts and returns a bearer token from the request.
// Returns an empty string if there is no bearer token or the token is invalid.
func (f *Factory_t) tokenFromBearerToken(r *http.Request) string {
	log.Printf("jot: bearer: entered\n")
	headerAuthText := r.Header.Get("Authorization")
	if headerAuthText == "" {
		return ""
	}
	log.Printf("jot: bearer: found authorization header\n")
	authTokens := strings.SplitN(headerAuthText, " ", 2)
	if len(authTokens) != 2 {
		return ""
	}
	log.Printf("jot: bearer: found authorization token\n")
	authType, authToken := authTokens[0], strings.TrimSpace(authTokens[1])
	if authType != "Bearer" {
		return ""
	}
	log.Printf("jot: bearer: found bearer token\n")
	return authToken
}

// tokenFromCookie extracts and returns a token from a cookie in the request.
// Returns an empty string if there is no cookie or the token is invalid.
func (f *Factory_t) tokenFromCookie(r *http.Request) string {
	log.Printf("jot: cookie: entered\n")
	c, err := r.Cookie(f.cookie.name)
	if err != nil {
		log.Printf("jot: cookie: %+v\n", err)
		return ""
	}
	log.Printf("jot: cookie: token\n")
	return c.Value
}

// Signer_i interface
type Signer_i interface {
	// Algorithm returns the name of the algorithm used by the signer.
	// The Factory will set the JOT header's "alg" field to this value when it is signed.
	// Example: "HS256"
	Algorithm() string

	// Expire clears the signer's expiration date.
	Expire()

	// Expired returns true if the Signer is expired.
	Expired() bool

	// Id is the unique identifier for this signer
	Id() string

	// Sign returns a slice containing the signature of the message.
	Sign(msg []byte) ([]byte, error)

	// Signed returns true if the msg was signed by this Signer.
	Signed(msg, signature []byte) bool
}

// HS256Signer_t implements a Signer using HMAC256.
type HS256Signer_t struct {
	id  string
	exp time.Time
	key []byte
}

func NewHS256Signer(id string, secret []byte, ttl time.Duration) (*HS256Signer_t, error) {
	return &HS256Signer_t{
		id:  id,
		key: append([]byte{}, secret...),
		exp: time.Now().Add(ttl).UTC(),
	}, nil
}

// Algorithm implements the Signer interface
func (s *HS256Signer_t) Algorithm() string {
	return "HS256"
}

// Expire implements the Signer interface.
func (s *HS256Signer_t) Expire() {
	s.exp = time.Unix(0, 0)
}

// Expired implements the Signer interface.
func (s *HS256Signer_t) Expired() bool {
	return s.exp.Before(time.Now().UTC())
}

// Id implements the Signer interface.
func (s *HS256Signer_t) Id() string {
	return s.id
}

// Sign implements the Signer interface.
func (s *HS256Signer_t) Sign(msg []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, s.key)
	if _, err := hm.Write(msg); err != nil {
		return nil, err
	}
	return hm.Sum(nil), nil
}

// Signed implements the Signer interface.
func (s *HS256Signer_t) Signed(msg, signature []byte) bool {
	ours, err := s.Sign([]byte(msg))
	return err == nil && bytes.Equal(signature, ours)
}

type Factory_t struct {
	sync.Mutex
	cookie struct {
		name string
		path string
	}
	ttl     time.Duration
	signers map[string]Signer_i
}

// AddSigner adds a new Signer to the pool
func (f *Factory_t) AddSigner(signer Signer_i) error {
	if signer.Expired() {
		return ErrSignerExpired
	}
	f.Lock()
	defer f.Unlock()
	f.signers[signer.Id()] = signer

	return nil
}

// DeleteSigner removes an existing Signer from the pool.
func (f *Factory_t) DeleteSigner(id string) {
	f.Lock()
	defer f.Unlock()
	delete(f.signers, id)
}

// DeleteExpiredSigners removes all expired Signers from the pool.
func (f *Factory_t) DeleteExpiredSigners() {
	f.Lock()
	defer f.Unlock()

	for id, signer := range f.signers {
		if signer.Expired() {
			delete(f.signers, id)
		}
	}
}

func (f *Factory_t) LookupSigner(id, algorithm string) (Signer_i, bool) {
	f.Lock()
	defer f.Unlock()

	signer, ok := f.signers[id]
	if !ok {
		return nil, false
	} else if signer.Expired() {
		delete(f.signers, id)
		return nil, false
	} else if signer.Algorithm() != algorithm {
		return nil, false
	}
	return signer, ok
}

var (
	// authenticatedUser is an empty user structure provided for convenience.
	unauthenticatedUser = User_t{Roles: map[string]bool{"authenticated": false}}
)

// payloadFromRequest returns the user data from the token in the request.
// If the token is valid (signed and not expired), we return the user data.
// Otherwise, we return an unauthenticated user.
func (f *Factory_t) payloadFromRequest(r *http.Request) (User_t, bool) {
	token := f.tokenFromRequest(r)
	if token == "" {
		return unauthenticatedUser, false
	}
	claim, err := f.claimsFromToken(token)
	if err != nil {
		return unauthenticatedUser, false
	}
	// return a copy of the payload so that we can release the claim's memory
	user := claim.Payload
	if user.Roles == nil {
		user.Roles = map[string]bool{}
	}
	user.Roles["authenticated"] = true
	return user, true
}

// claimsFromToken extracts claims from a token.
// It returns an error if the token is invalid, expired, or hasn't been signed correctly.
// Otherwise, it returns the claims from the token.
func (f *Factory_t) claimsFromToken(token string) (*Claims_t, error) {
	// extract the header, claims, and signature from the token
	fields := strings.Split(token, ".")
	if len(fields) != 3 {
		return nil, ErrInvalidToken
	} else if len(fields[0]) > 99 {
		// header should be about 10 bytes for `typ`, 12 for `alg` and 40 for `kid`
		return nil, ErrInvalidHeader
	}
	// assign the fields to header, claims, and signature
	h64, c64, s64 := fields[0], fields[1], fields[2]

	// decode the header
	var header Header_t
	if data, err := decode_str(h64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &header); err != nil {
		return nil, err
	} else if header.TokenType != "JOT" {
		return nil, ErrUnknownType
	}

	// decode the signature
	signature, err := decode_str(s64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// use the header to find the original signer and confirm the signature
	signer, ok := f.signers[header.KeyID]
	if !ok {
		return nil, ErrInvalidSigner
	} else if signer.Algorithm() != header.Algorithm {
		return nil, ErrInvalidAlgorithm
	} else if !signer.Signed([]byte(h64+"."+c64), signature) {
		return nil, ErrInvalidSignature
	}

	// decode and validate the claims
	var claims Claims_t
	if data, err := decode_str(c64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &claims); err != nil {
		return nil, err
	} else if !claims.IsNotExpired(time.Now().UTC()) {
		return nil, ErrClaimsExpired
	}

	// return the claims if the message is signed and they haven't expired
	return &claims, nil
}

func (f *Factory_t) NewTokenCookie(ttl time.Duration, payload User_t) (*http.Cookie, error) {
	signer, err := f.getSigner()
	if err != nil {
		return nil, ErrMissingSigner
	}

	// update and marshal the header to JSON
	h64, err := Header_t{
		Algorithm: signer.Algorithm(),
		KeyID:     signer.Id(),
		TokenType: "JOT",
	}.Encode()
	if err != nil {
		return nil, err
	}

	// update and marshal the claims to JSON
	iat := time.Now().UTC()
	exp := iat.Add(ttl).UTC()
	c64, err := Claims_t{
		IssuedAt:  NumericDate_t(iat),
		ExpiresAt: NumericDate_t(exp),
		Payload:   payload,
	}.Encode()
	if err != nil {
		return nil, err
	}

	// create message as header + '.' + claims
	var msg []byte
	msg = append(msg, h64...)
	msg = append(msg, '.')
	msg = append(msg, c64...)

	// sign the message
	signature, err := signer.Sign(msg)
	if err != nil {
		return nil, err
	}
	s64 := encode_bytes(signature)

	// return the token as message + '.' + signature
	msg = append(msg, '.')
	msg = append(msg, s64...)
	return &http.Cookie{
		Path:     f.cookie.path, // Path for which the cookie is valid
		Name:     f.cookie.name,
		Value:    string(msg),
		Expires:  exp,
		HttpOnly: true, // make sure HttpOnly is true to prevent javascript access
		Secure:   true, // make sure Secure is true if over HTTPS
	}, nil
}

// getSigner returns a signer from the pool.
// The pool is a map, so the signer might change between runs.
// Warning: has the side-effect of deleting expired signers from the pool.
func (f *Factory_t) getSigner() (Signer_i, error) {
	f.Lock()
	defer f.Unlock()

	for id, signer := range f.signers {
		if !signer.Expired() {
			return signer, nil
		}
		delete(f.signers, id)
	}
	return nil, ErrNotFound
}

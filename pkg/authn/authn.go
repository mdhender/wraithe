/*
 * wraith - a game engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package authn implements an authentication store and service.
// Note: we use "user" for both end-users, clients, and other services.
package authn

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

// Store is our in-memory store for identification records.
type Store struct {
	// identities stores all identification records.
	// The map key is UserId.
	identities map[string]*Identity

	// keys stores all signing keys.
	// The map key is SigningKey.KeyId
	keys map[string]*SigningKey

	// ttl is the maximum time-to-live for tokens.
	// it is constrained by the expiration of the signing key.
	ttl time.Duration
}

func New(filename string, ttl time.Duration) (*Store, error) {
	return &Store{
		identities: make(map[string]*Identity),
		keys:       make(map[string]*SigningKey),
		ttl:        ttl,
	}, nil
}

func (s *Store) Authenticate(userId, userSecret string) Token {
	if s == nil || s.identities == nil {
		return Token{}
	}
	id, ok := s.identities[userId]
	if !ok {
		return Token{}
	}

	hashedSecret := []byte("do-a-bcrypt-here")
	if !bytes.Equal(hashedSecret, id.HashedSecret) {
		return Token{}
	}

	t := Token{
		Id:        id.Id,
		ExpiresAt: time.Now().Add(s.ttl).UTC(),
	}

	// find a valid signing key
	for _, key := range s.keys {
		// never use a key that will expire before the token
		if t.ExpiresAt.Before(key.expiresAt) {
			continue
		}
		// use the key to sign the token
		if err := t.Sign(key); err != nil {
			continue
		}
		// we have signed our token, so quit the loop
		break
	}

	// note: it is possible that the store has no valid keys.
	// if that happens, our token will not be signed and will not be valid.
	return t
}

// Identity is used to store a user's information for authentication.
type Identity struct {
	Id           string `json:"id"`
	UserId       string `json:"user-id"`
	HashedSecret []byte `json:"hashed-secret,omitempty"`
}

// Token is returned from an authentication challenge.
type Token struct {
	Id        string    // Identity.Id
	ExpiresAt time.Time // moment (in UTC) the token expires
	KeyId     string    // SigningKey.KeyId
	Signature []byte    // generated signature
}

func (t *Token) Sign(k *SigningKey) (err error) {
	t.KeyId = k.keyId
	var msg string
	t.Signature, err = k.Sign(msg)
	if err != nil {
		return err
	}
	return nil
}

type Authenticator struct {
	ttl time.Duration
}

type Authentication struct {
	ID            string    // the id of the user
	expiresAt     time.Time // the moment this authentication becomes invalid
	authenticated bool
}

// IsAuthenticated returns true only if the user has authenticated successfully.
func (a Authentication) IsAuthenticated() bool {
	return a.authenticated
}

// IsValid returns true only if the user has been authenticated and not expired.
func (a Authentication) IsValid() bool {
	return a.IsAuthenticated() && time.Now().Before(a.expiresAt)
}

// FromRequest extracts credentials from the http request.
// It adds the result to the request's context.
func (a *Authenticator) FromRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var an Authentication

		// source for authentication depends on the type of request.
		switch r.Header.Get("Content-type") {
		case "application/json":
			// request is RESTish, so try to fetch a bearer token
		case "text/html":
			// request is plain HTML, so try to fetch a cookie
		default:
			// don't authenticate other content types
		}

		// add the authentication to the context
		ctx := context.WithValue(r.Context(), "authn", an)

		// call the next handler in the chain,
		// passing the response writer and the updated request object with the new context value.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func fromCookie(r http.Request) Authentication {
	return Authentication{}
}

func fromToken(r http.Request) Authentication {
	return Authentication{}
}

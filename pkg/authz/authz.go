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

// Package authz implements an authorization middleware.
// It relies on the authn package adding an Authorization to the request context.
package authz

import (
	"context"
	"net/http"
)

type Authorizer struct{}

type Authorization struct{}

func (a *Authorizer) FromRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var az Authorization

		// fetch the authorization from the context

		// if the request is HTML (not RESTish), try to fetch a cookie

		// add the authentication to the context
		ctx := context.WithValue(r.Context(), "authz", az)

		// call the next handler in the chain,
		// passing the response writer and the updated request object with the new context value.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

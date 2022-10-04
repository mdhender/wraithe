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

// Package main implements the shell for wraith.
package main

import (
	"context"
	"errors"
	"github.com/mdhender/wraithe/pkg/cedar"
	"github.com/mdhender/wraithe/pkg/cfg"
	"github.com/mdhender/wraithe/pkg/cli"
	"github.com/mitchellh/go-homedir"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	defer func(started time.Time) {
		elapsed := time.Now().Sub(started)
		log.Printf("wraith: total time %v\n", elapsed)
	}(time.Now())

	// find starting directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// find home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%-30s == %q\n", "home", home)

	// seed the default PRNG source.
	seed, err := cedar.Seed()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%-30s == %d\n", "seed", seed)
	rand.Seed(seed)

	if err := cli.Execute(nil); err != nil {
		log.Fatal(err)
	}

	c, err := cfg.Read(filepath.Join("testdata", "wraith.cfg"))
	if err != nil {
		log.Fatal(err)
	}
	c.Home = home
	c.PRNG.Seed = seed
	c.WorkingDir = cwd

	err = run(c)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cfg.Config) error {
	if c == nil {
		return cli.Execute(c)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	s := &server{}
	s.Addr = ":8080"
	s.BaseContext = func(l net.Listener) context.Context {
		ctx = context.WithValue(ctx, "addr", l.Addr().String())
		return ctx
	}
	s.Handler = s.routes()

	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	cancelCtx()

	return err

}

type server struct {
	http.Server
}

// shiftPath splits the given path into the first segment (head) and
// the rest (tail). For example, "/foo/bar/baz" gives "foo", "/bar/baz".
func shiftPath(p string) (head, tail string) {
	// from https://benhoyt.com/writings/go-routing/#shiftpath
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// ensureMethod is a helper that reports whether the request's method is
// the given method, writing an Allow header and a 405 Method Not Allowed
// if not. The caller should return from the handler if this returns false.
func ensureMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	// from https://benhoyt.com/writings/go-routing/#shiftpath
	if method != r.Method {
		w.Header().Set("Allow", method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

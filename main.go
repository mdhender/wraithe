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
	"github.com/mdhender/wraithe/pkg/cfg"
	"github.com/mdhender/wraithe/pkg/cli"
	"log"
	"path/filepath"
	"time"
)

func main() {
	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	defer func(started time.Time) {
		elapsed := time.Now().Sub(started)
		log.Printf("wraith: total time %v\n", elapsed)
	}(time.Now())

	c, err := cfg.Read(filepath.Join("testdata", "wraith.cfg"))
	if err != nil {
		log.Fatal(err)
	}

	err = run(c)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cfg.Config) error {
	return cli.Execute(c)
}

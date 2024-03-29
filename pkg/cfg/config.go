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

// Package cfg implements a store for configuration data.
package cfg

// Config is the configuration data.
// All fields are public.
type Config struct {
	Meta       Meta   `json:"meta"`
	Home       string `json:"home"`
	WorkingDir string `json:"working-dir"`
	Server     Server `json:"server"`
	PRNG       PRNG   `json:"prng"`
}

// Meta is meta-data about the configuration file.
type Meta struct {
	FileName string `json:"file-name,omitempty"`
}

// PRNG is configuration for random numbers
type PRNG struct {
	Seed int64
}

// Server is the configuration for the wraith server.
type Server struct {
	Http Http `json:"http"`
}

// Http is the configuration for the http server.
type Http struct {
	Port string `json:"port,omitempty"`
}

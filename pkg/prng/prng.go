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

// Package prng implements a motley assortment of psuedo-random number generators.
package prng

// PRNG is a generator.
type PRNG func() uint32

func (p PRNG) Roll(n, d int) (result int) {
	if n < 1 || d < 1 {
		return 0
	}

	for ; n > 0; n-- {
		result += int(p() % uint32(d))
	}

	return result
}

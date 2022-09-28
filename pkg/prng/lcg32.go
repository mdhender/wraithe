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

package prng

// LCG32 is from open adventure, maybe?
func LCG32(x uint32) func() uint32 {
	x = x % 1048576

	// generate and return the next value.
	fn := func() uint32 {
		x = (x*1093 + 221587) % 1048576
		return x<<21 | x>>11
	}

	// source recommends tossing the first output
	for i := 0; i < 1; i++ {
		fn()
	}

	return fn
}

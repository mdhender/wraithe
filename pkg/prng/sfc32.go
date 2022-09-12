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

// SFC32 is from https://simblob.blogspot.com/2022/05/upgrading-prng.html#more
//
// Example
//   sfc32 := prng.SFC32(0, 12345, 0, 1)
//   for i := 0; i < 10; i++ {
//   	log.Println(sfc32())
//   }
//
// Expected Output:
//   235160590
//   2967261163
//   116171463
//   2882324903
//   362604721
//   4227106926
//   1933307004
//   1608300071
//   2256615412
//   2701957640
//
func SFC32(a, b, c, d uint32) func() uint32 {
	fn := func() uint32 {
		t := (a + b) + d
		d = d + 1
		a = b ^ b>>9
		b = c + (c << 3)
		c = c<<21 | c>>11
		c = c + t
		return t
	}

	// source recommends tossing first 12 outputs
	for i := 0; i < 12; i++ {
		fn()
	}

	return fn
}

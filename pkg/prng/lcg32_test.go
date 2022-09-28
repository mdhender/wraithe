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

import "testing"

func TestLCG(t *testing.T) {
	lcg32 := LCG32(0)
	for _, expect := range []uint32{1178599519, 564134195, 3263168954, 2665480396, 2227175438, 4196401256, 486539424, 56623112, 2604662946, 178258093} {
		if got := lcg32(); got != expect {
			t.Errorf("wanted %d: got %d\n", expect, got)
		}
	}

	lcg32 = LCG32(1)
	for _, expect := range []uint32{2573205670, 155189689, 2963276144, 1300234382, 367002109, 2545943034, 421527894, 2013266181, 2321547461, 4265607635} {
		if got := lcg32(); got != expect {
			t.Errorf("wanted %d: got %d\n", expect, got)
		}
	}
}

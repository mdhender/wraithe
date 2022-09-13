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

import (
	"testing"
)

func TestSFC32(t *testing.T) {
	sfc32 := SFC32(0, 12345, 0, 1)
	for _, expect := range []uint32{
		235160590,
		2967261163,
		116171463,
		2882324903,
		362604721,
		4227106926,
		1933307004,
		1608300071,
		2256615412,
		2701957640,
	} {
		if got := sfc32(); got != expect {
			t.Errorf("wanted %d: got %d\n", expect, got)
		}
	}
}

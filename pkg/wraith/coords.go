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

package wraith

import (
	"fmt"
	"math"
)

// coords are the location of a unit in the game.
type coords struct {
	x, y, z float64
}

// distance returns the distance between two points
func (c coords) distance(a coords) float64 {
	dx, dy, dz := c.x-a.x, c.y-a.y, c.z-a.z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// roundToInt returns the coordinates rounded off
func (c coords) roundToInt() coords {
	return coords{x: math.Round(c.x), y: math.Round(c.y), z: math.Round(c.z)}
}

// xyz returns a string with the rounded coordinates
func (c coords) xyz() string {
	return fmt.Sprintf("%3d%3d%3d", int(math.Round(c.x)), int(math.Round(c.y)), int(math.Round(c.z)))
}

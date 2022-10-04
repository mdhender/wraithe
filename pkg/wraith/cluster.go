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
	"bytes"
	"fmt"
	"github.com/mdhender/wraithe/pkg/prng"
	"html/template"
	"log"
	"math"
	"math/rand"
	"path/filepath"
)

// note: we use the term "ring" to include the stars and systems that are
//       a given distance from the origin of the cluster. ring zero is the
//       origin, ring 1 is 1 light year from it, etc.

var (
	// systems are not uniformly distributed in the cluster.
	// instead they are generated based on the ring they're in.
	//expectedSystemsPerRing [16]int = [16]int{1, 4, 6, 8, 9, 11, 13, 11, 10, 9, 8, 8, 7, 6, 5, 5}
	systemsPerRing = [16]int{1, 6, 8, 8, 9, 9, 8, 6, 6, 6, 6, 6, 4, 4, 3, 3}
)

type cluster struct {
	prng    *prng.PRNG
	systems []*system
}
type system struct {
	ring   int
	coords coords
	stars  []*star
}
type star struct {
	system *system
}
type coords struct {
	x, y, z float64
}

// a cluster has a radius of about 15.

// F generates a new cluster.
// `n` is the number of players (nations) to populate.
func F(minStars int) *cluster {
	c := cluster{}

	totalSystems := 0
	for _, n := range systemsPerRing {
		totalSystems += n
	}
	log.Printf("[f] generating cluster with %3d systems and %3d stars\n", totalSystems, minStars)

	// the "home world" is always in the "home system" which is always at (0,0,0).
	// it always contains 1 star and 10 planets.
	// the resources on the home world are determined by the number of players in the game.
	c.systems = append(c.systems, &system{ring: 0, coords: coords{0, 0, 0}})

	// the first ring always contains six systems, one at each cardinal point.
	c.systems = append(c.systems, &system{ring: 1, coords: coords{1, 0, 0}})
	c.systems = append(c.systems, &system{ring: 1, coords: coords{-1, 0, 0}})
	c.systems = append(c.systems, &system{ring: 1, coords: coords{0, 1, 0}})
	c.systems = append(c.systems, &system{ring: 1, coords: coords{0, -1, 0}})
	c.systems = append(c.systems, &system{ring: 1, coords: coords{0, 0, 1}})
	c.systems = append(c.systems, &system{ring: 1, coords: coords{0, 0, -1}})

	for ring, expectedSystems := range systemsPerRing {
		if ring == 0 || ring == 1 {
			continue
		}

		log.Printf("[f] generating ring %2d with %2d systems\n", ring, expectedSystems)

		for i := 0; i < expectedSystems; i++ {
			sys := &system{ring: ring}

			// create a location for the system by generating 15 random
			// points in the ring and using the one that is the furthest
			// from all existing systems
			maxDistance := 0.0 // maximum distance of points so far
			for pn := 0; pn < 15; pn++ {
				if pt := randomPoint(ring); pn == 0 || maxDistance < minDistance(pt, c.systems) {
					sys.coords = pt
				}
			}
			sys.coords.x = math.Round(sys.coords.x)
			sys.coords.y = math.Round(sys.coords.y)
			sys.coords.z = math.Round(sys.coords.z)
			log.Println(ring, sys.coords)

			c.systems = append(c.systems, sys)
		}
	}

	// create one star in every system
	for _, sys := range c.systems {
		sys.stars = append(sys.stars, &star{system: sys})
	}

	// add any remaining stars to random systems.
	for remaining := minStars - len(c.systems); remaining > 0; remaining-- {
		// pick a system at random. well, sort of random. we don't want binary
		// systems in rings 0 or 1, so if the pick lands there, pick again.
		sys := c.systems[rand.Intn(len(c.systems))]
		for sys.ring == 0 || sys.ring == 1 {
			sys = c.systems[rand.Intn(len(c.systems))]
		}
		// and add a star to it
		sys.stars = append(sys.stars, &star{system: sys})
	}

	totalStars := 0
	for _, sys := range c.systems {
		totalStars += len(sys.stars)
		log.Printf("[ring] %2d: %5.0f,%5.0f,%5.0f %2d stars %4d\n", sys.ring, sys.coords.x, sys.coords.y, sys.coords.z, len(sys.stars), totalStars)
	}

	return &c
}

// G generates a new cluster.
// The cluster will have 128 stars in 256 systems.
func G(minSystems, minStars int, scale float64) *cluster {
	const probes = 15

	starsLeft := minStars

	c := cluster{
		systems: []*system{
			//{coords: coords{x: 0, y: -15, z: 0}}, // north pole
			{coords: coords{x: 0, y: 0, z: 0}}, // center
			//{coords: coords{x: 0, y: 15, z: 0}},  // south pole
		},
	}

	// create the systems.
	// we will add stars to systems until the number left is zero.
	// the remaining systems will have no stars.
	// systems with stars can't be closer than 4 light years to another system with stars.
	// systems without stars can't be closer than 1 light year to another system.
	for len(c.systems) < minSystems {
		sys := &system{}
		if starsLeft > 0 {
			sys.stars, starsLeft = append(sys.stars, &star{}), starsLeft-1
			if starsLeft > 0 && len(c.systems) < 28 {
				sys.stars, starsLeft = append(sys.stars, &star{}), starsLeft-1
				if starsLeft > 0 && len(c.systems) < 12 {
					sys.stars, starsLeft = append(sys.stars, &star{}), starsLeft-1
					if starsLeft > 0 && len(c.systems) < 6 {
						sys.stars, starsLeft = append(sys.stars, &star{}), starsLeft-1
						if starsLeft > 0 && len(c.systems) < 3 {
							sys.stars, starsLeft = append(sys.stars, &star{}), starsLeft-1
						}
					}
				}
			}
		}

		// loop until we get a point that isn't near the center or too close to an existing system
		var ring int
		var spt coords
		maxDistance, closestNeighbor := 0.0, float64(2+len(sys.stars)*2)
		if len(sys.stars) == 0 {
			closestNeighbor = 1.0
		}
		for probe := 0; probe < probes; {
			pt := getPoint(scale)
			if ring = int(math.Round(pt.distance(coords{}))); ring < 5 || minDistance(pt, c.systems) < closestNeighbor {
				continue
			}
			if d := minDistance(pt, c.systems); d > maxDistance {
				spt, maxDistance = pt, d
			}
			probe++
		}
		sys.ring = int(math.Round(spt.distance(coords{})))
		sys.coords = spt.roundToInt()
		//if starsLeft > 0 {
		//	log.Printf("[ring] %2d: %7.2f,%7.2f,%7.2f stars: %5d / %5d\n", sys.ring, sys.coords.x, sys.coords.y, sys.coords.z, len(sys.stars), starsLeft)
		//}
		c.systems = append(c.systems, sys)
	}

	return &c
}

func getPoint(scale float64) coords {
	u, v, d := rand.Float64(), rand.Float64(), rand.Float64()
	theta, phi, r := u*2.0*math.Pi, math.Acos(2.0*v-1.0), d //math.Cbrt(d)
	sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)
	sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)
	return coords{x: scale * r * sinPhi * cosTheta, y: scale * r * sinPhi * sinTheta, z: scale * r * cosPhi}
}

// minDistance returns the distance from the coordinates to the nearest system.
func minDistance(c coords, systems []*system) (d float64) {
	for i := range systems {
		if ds := c.distance(systems[i].coords); i == 0 || ds < d {
			d = ds
		}
	}
	return d
}

// randomPoint returns a random location that is the given distance from the origin.
// uses the method from https://www.cs.cmu.edu/~mws/rpos.html
//
//	    z is in range -R ... R
//	  phi is in range  0 ... 2pi
//	theta is sin-1(z/R)
//	    x is R cos(theta) cos(phi)
//	    y is R cos(theta) sin(phi)
func randomPoint(distance int) (c coords) {
	R := float64(distance)
	c = coords{z: R * rand.Float64()}
	if rand.Int()%2 == 0 {
		c.z = -c.z
	}
	phi := 2 * math.Pi * rand.Float64()
	rCosTheta := R * math.Cos(math.Asin(c.z/R))
	c.x, c.y = rCosTheta*math.Cos(phi), rCosTheta*math.Sin(phi)
	return c
}

func (c coords) distance(a coords) float64 {
	dx, dy, dz := c.x-a.x, c.y-a.y, c.z-a.z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c coords) roundToInt() coords {
	return coords{x: math.Round(c.x), y: math.Round(c.y), z: math.Round(c.z)}
}

func (c coords) xyz() string {
	return fmt.Sprintf("%3d%3d%3d", int(math.Round(c.x)), int(math.Round(c.y)), int(math.Round(c.z)))
}

func (c *cluster) ToHTML(templates string, tname string, tfm template.FuncMap) ([]byte, error) {
	type System struct {
		Ring    int
		Size    int
		Color   string
		X, Y, Z float64
	}
	var data struct {
		Systems []*System
	}

	for _, sys := range c.systems {
		var color string
		switch len(sys.stars) {
		case 5:
			color = "red"
		case 4:
			color = "teal"
		case 3:
			color = "blue"
		case 2:
			color = "green"
		case 1:
			color = "silver"
		default:
			color = "grey"
		}
		data.Systems = append(data.Systems, &System{Ring: sys.ring, X: sys.coords.x, Y: sys.coords.y, Z: sys.coords.z, Size: len(sys.stars), Color: color})
	}
	log.Printf("[cluster] len(data.Systems) is %d\n", len(data.Systems))

	t, err := template.New(tname).Funcs(tfm).ParseFiles(filepath.Join(templates, tname))
	if err != nil {
		return nil, err
	}
	bw := &bytes.Buffer{}
	if err = t.Execute(bw, data); err != nil {
		return nil, err
	}

	return bw.Bytes(), nil
}

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
	"github.com/mdhender/wraithe/pkg/prng"
	"html/template"
	"log"
	"path/filepath"
)

// cluster is a container for the systems
type cluster struct {
	prng    *prng.PRNG
	systems []*system
}

// ToHTML returns a pretty picture of the cluster.
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

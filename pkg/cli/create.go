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

package cli

import (
	"github.com/mdhender/wraithe/pkg/wraith"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"os"
)

var createArgs struct {
	minStars int
}

// createCmd implements the commands needed to create a new game.
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create a new game",
	Long:    `Create a new game.`,
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		//wraith.F(createArgs.minStars)
		c := wraith.G(512, 128, 15.0)
		b, err := c.ToHTML("D:\\wraithe\\templates", "cluster.gohtml", template.FuncMap{})
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		fname := "D:\\wraithe\\testdata\\gencluster.html"
		err = os.WriteFile(fname, b, 0666)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		log.Printf("[create] created %q\n", fname)
	},
}

func init() {
	cmdCLI.AddCommand(createCmd)
	createCmd.Flags().IntVar(&createArgs.minStars, "min-stars", 125, "minimum number of stars in the cluster")
}

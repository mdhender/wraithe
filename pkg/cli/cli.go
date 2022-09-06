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

// Package cli implements the command line interface for wraith.
package cli

import (
	"github.com/mdhender/wraithe/pkg/cedar"
	"github.com/mdhender/wraithe/pkg/cfg"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
)

// cmdCLI represents the base command when called without any subcommands
var cmdCLI = &cobra.Command{
	Use:   "wraith",
	Short: "Wraith game engine",
	Long: `wraith is the game engine for Wraith.
This application provides an API to the game engine.`,
	Version: "0.0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// find starting directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		log.Printf("%-30s == %q\n", "cwd", cwd)

		// find home directory
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		log.Printf("%-30s == %q\n", "home", home)

		// seed the default PRNG source.
		seed, err := cedar.Seed()
		if err != nil {
			return err
		}
		log.Printf("%-30s == %d\n", "seed", seed)
		rand.Seed(seed)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// do nothing for the base command?
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the root Command.
func Execute(c *cfg.Config) error {
	return cmdCLI.Execute()
}

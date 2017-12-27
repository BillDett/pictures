// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pictures/db"
	"pictures/index"
)

var from string
var idx string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an index or database file.",
	Long:  `When building an index, the root path of pictures is provided.  When building a database a database filename is provided.  If the database file does not exist it will be created.  If it already exists it will be updated.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("build command requires argument 'database' or 'index'")
		}
		if args[0] != "database" && args[0] != "index" {
			return errors.New("Invalid argument " + args[0] + " given.  Must be 'database' or 'index'")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if from == "" {
			fmt.Fprintln(os.Stderr, "Error: must provide a value for --from")
			return
		}
		if args[0] == "index" {
			index.Photoindex(from)
		} else if args[0] == "database" {
			if idx == "" {
				fmt.Fprintln(os.Stderr, "Error: must provide a value for --index when building a database")
				return
			}
			db.Photodatabase(from, idx)
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.PersistentFlags().StringVar(&from, "from", "", "Source file/directory from which to start build.")
	buildCmd.PersistentFlags().StringVar(&idx, "index", "", "Index file to be used when building a database.")
}

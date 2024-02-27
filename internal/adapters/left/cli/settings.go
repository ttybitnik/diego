/*
   DIEGO - A data importer extension for Hugo
   Copyright (C) 2024 Vinicius Moraes <vinicius.moraes@eternodevir.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cli

import (
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show current settings",
	Long:  `Show current settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Error reading config file:", err)
		}

		settings := viper.AllSettings()
		printSettings("", settings)

		if settings == nil {
			log.Fatalln("Error: no configuration file found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}

func printSettings(prefix string, settings map[string]interface{}) {
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := settings[key]
		if nestedMap, ok := value.(map[string]interface{}); ok {
			printSettings(fmt.Sprintf("%s%s.", prefix, key), nestedMap)
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}
	}
}

/*
   DIEGO - A data importer extension for Hugo
   Copyright (C) 2024 Vinícius Moraes <vinicius.moraes@eternodevir.com>

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultsCmd = &cobra.Command{
	Use:  "defaults",
	Args: cobra.MatchAll(cobra.ExactArgs(0)),
	Example: "diego set defaults true\n" +
		"diego s defaults 1",
	Short: "Restore Diego default settings",
	Long:  "Restore Diego default settings.",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("diego.import.all", false)
		viper.Set("diego.import.overwrite", false)
		viper.Set("diego.import.scrape", false)
		viper.Set("diego.import.shortcode", false)
		viper.Set("diego.import.format", "yaml")
		viper.Set("diego.import.hugodir", ".")

		fmt.Println("Diego default settings restored.")

		err := viper.WriteConfig()
		if err != nil {
			log.Fatalln("Error writing into config file:", err)
			return
		}
	},
}

func init() {
	setCmd.AddCommand(defaultsCmd)
}

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

var shortcodeCmd = &cobra.Command{
	Use:       "shortcode {true|false|1|0|enabled|disabled}",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"true", "false", "1", "0", "enabled", "disabled"},
	Example: "diego set shortcode true\n" +
		"diego s shortcode false",
	Short: "Enable or disable the shortcode flag by default",
	Long:  "Enable or disable the shortcode flag by default.",
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "true", "1", "enabled":
			viper.Set("diego.import.shortcode", true)
			fmt.Println("Shortcode flag is set to 'enabled' by default.")
		case "false", "0", "disabled":
			viper.Set("diego.import.shortcode", false)
			fmt.Println("Shortcode flag is set to 'disabled' by default.")
		}

		err := viper.WriteConfig()
		if err != nil {
			log.Fatalln("Error writing into config file:", err)
			return
		}
	},
}

func init() {
	setCmd.AddCommand(shortcodeCmd)
}

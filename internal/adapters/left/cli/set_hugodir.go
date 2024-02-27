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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hugodirCmd = &cobra.Command{
	Use:  "hugodir path",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Example: "diego set hugodir ~/Projects/HugoBlog \n" +
		"diego s hugodir ../Projects/HugoSite",
	Short: "Set path to the Hugo directory (default \".\")",
	Long:  "Set path to the Hugo directory (default \".\").",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := args[0]

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			log.Fatalln("Invalid path "+dirPath+":", err)
		}

		dirPath = filepath.Clean(dirPath)

		viper.Set("diego.import.hugodir", dirPath)
		fmt.Printf("Format flag is set to '%s' by default.\n", dirPath)

		err := viper.WriteConfig()
		if err != nil {
			log.Fatalln("Error writing into config file:", err)
			return
		}
	},
}

func init() {
	setCmd.AddCommand(hugodirCmd)
}

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
	"github.com/ttybitnik/diego/internal/app/domain"

	"github.com/spf13/cobra"
)

var instapaperCmd = &cobra.Command{
	Use:       "instapaper {list} -i file",
	Aliases:   []string{"ip"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"list"},
	Example: "diego import instapaper list -i list.csv\n" +
		"diego i ip list -i list.csv --all --shortcode",
	Short: "Import data from Instapaper",
	Long:  `Import data from Instapaper.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		ca.cobra.input = input

		switch args[0] {
		case "list":
			ca.cobra.model = domain.InstapaperList
			ca.cobraImport(ca.cobra)
		}
	},
}

func init() {
	importCmd.AddCommand(instapaperCmd)
	ca = cobraAdapter()

	instapaperCmd.Flags().BoolVar(&ca.cobra.all, "all", false, "import every available field from CSV file")
	instapaperCmd.Flags().BoolVar(&ca.cobra.overwrite, "overwrite", false, "overwrite existent output data file")
	instapaperCmd.Flags().BoolVar(&ca.cobra.shortcode, "shortcode", false, "generate a shortcode template for Hugo")
	instapaperCmd.Flags().StringVar(&ca.cobra.format, "format", domain.OutputYAML, "output format for the Hugo data file")
	instapaperCmd.Flags().StringVar(&ca.cobra.hugodir, "hugodir", ".", "path to the Hugo directory")
	instapaperCmd.Flags().StringVarP(&ca.cobra.input, "input", "i", "", "path to the CSV file (required)")
}

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
	"github.com/ttybitnik/diego/internal/app/domain"

	"github.com/spf13/cobra"
)

var youtubeCmd = &cobra.Command{
	Use:       "youtube {playlist|subscriptions} -i file",
	Aliases:   []string{"y"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"playlist", "subscriptions"},
	Example: "diego import youtube subscriptions -i subscriptions.csv\n" +
		"diego i y playlist -i playlist.csv --all --scrape --shortcode",
	Short: "Import data from YouTube",
	Long:  `Import data from YouTube.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		ca.cobra.input = input

		switch args[0] {
		case "playlist":
			ca.cobra.model = domain.YoutubePlaylist
			ca.cobraImport(ca.cobra)
		case "subscriptions":
			ca.cobra.model = domain.YoutubeSubscriptions
			ca.cobraImport(ca.cobra)
		}
	},
}

func init() {
	importCmd.AddCommand(youtubeCmd)
	ca = cobraAdapter()

	youtubeCmd.Flags().BoolVar(&ca.cobra.all, "all", false, "import every available field from CSV file")
	youtubeCmd.Flags().BoolVar(&ca.cobra.overwrite, "overwrite", false, "overwrite existent output data file")
	youtubeCmd.Flags().BoolVar(&ca.cobra.scrape, "scrape", false, "fetch additional data from CSV links using HTTP")
	youtubeCmd.Flags().BoolVar(&ca.cobra.shortcode, "shortcode", false, "generate a shortcode template for Hugo")
	youtubeCmd.Flags().StringVar(&ca.cobra.format, "format", domain.OutputYAML, "ouput format for the Hugo data file")
	youtubeCmd.Flags().StringVar(&ca.cobra.hugodir, "hugodir", ".", "path to the Hugo directory")
	youtubeCmd.Flags().StringVarP(&ca.cobra.input, "input", "i", "", "path to the CSV file (required)")
}

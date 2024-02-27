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
	"log"

	"github.com/ttybitnik/diego/internal/app/domain"

	"github.com/spf13/cobra"
)

var spotifyCmd = &cobra.Command{
	Use:       "spotify {library|playlist} -i file",
	Aliases:   []string{"s"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"library", "playlist", "watchlist"},
	Example: "diego import spotify library -i library.json\n" +
		"diego i s playlist -i playlist.json --all --shortcode",
	Short: "Import data from Spotify",
	Long:  `Import data from Spotify.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		ca.cobra.input = input

		switch args[0] {
		case "library":
			ca.cobra.model = domain.SpotifyLibrary
			ca.cobraImport(ca.cobra)
		case "playlist":
			ca.cobra.model = domain.SpotifyPlaylist
			ca.cobraImport(ca.cobra)
		}
	},
}

func init() {
	importCmd.AddCommand(spotifyCmd)
	ca = cobraAdapter()

	spotifyCmd.Flags().BoolVar(&ca.cobra.all, "all", false, "import every available field from JSON file")
	spotifyCmd.Flags().BoolVar(&ca.cobra.overwrite, "overwrite", false, "overwrite existent output data file")
	spotifyCmd.Flags().BoolVar(&ca.cobra.scrape, "scrape", false, "fetch additional data from JSON links using HTTP")
	spotifyCmd.Flags().BoolVar(&ca.cobra.shortcode, "shortcode", false, "generate a shortcode template for Hugo")
	spotifyCmd.Flags().StringVar(&ca.cobra.format, "format", domain.OutputYAML, "output format for the Hugo data file")
	spotifyCmd.Flags().StringVar(&ca.cobra.hugodir, "hugodir", ".", "path to the Hugo directory")
	spotifyCmd.Flags().StringVarP(&ca.cobra.input, "input", "i", "", "path to the JSON file (required)")

	err := spotifyCmd.Flags().MarkHidden("scrape")
	if err != nil {
		log.Fatalln("Error marking flag as hidden:", err)
	}
}

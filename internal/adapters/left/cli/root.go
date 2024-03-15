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

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	ciVersion string = "0.2.1" // x-release-please-version
)

var rootCmd = &cobra.Command{
	Use: "diego",
	Version: ciVersion + `
Copyright (C) 2024 Vinicius Moraes
License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>
This is free software, and you are welcome to change and redistribute it under
certain conditions. This program comes with ABSOLUTELY NO WARRANTY.`,
	Short: "DIEGO - A data importer extension for Hugo",
	Long: `DIEGO - A data importer extension for Hugo

Diego integrates with Hugo as a CLI tool to assist in importing and utilizing
exported social media data from various services on Hugo websites.

Complete documentation is available at <https://github.com/ttybitnik/diego>.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	if os.Getenv("DIEGO_GENDOCS") == "1" {
		generateDocs()
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/diego/config.yaml)")
}

func initConfig() {
	viper.SetDefault("diego.import.all", false)
	viper.SetDefault("diego.import.overwrite", false)
	viper.SetDefault("diego.import.scrape", false)
	viper.SetDefault("diego.import.shortcode", false)
	viper.SetDefault("diego.import.format", "yaml")
	viper.SetDefault("diego.import.hugodir", ".")

	// Use config file from the flag or file system.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home + "/.config/diego")

		_ = viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func generateDocs() {
	header := &doc.GenManHeader{
		Title:   "DIEGO",
		Section: "1",
		Manual:  "User Commands",
		Source:  "diego manual",
	}

	rootCmd.DisableAutoGenTag = true

	docsDir := ("docs")
	cobra.CheckErr(doc.GenManTree(rootCmd, header, docsDir+"/man"))
	cobra.CheckErr(doc.GenMarkdownTree(rootCmd, docsDir+"/help"))

	log.Println("Man and markdown docs recreated:", "./"+docsDir)
}

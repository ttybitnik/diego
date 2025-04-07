/*
   DIEGO - A data importer extension for Hugo
   Copyright (C) 2024, 2025 Vinícius Moraes <vinicius.moraes@eternodevir.com>

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

// Package cli contains the command line interface adapter logic.
package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ttybitnik/diego/internal/adapters/right/filesystem"
	"github.com/ttybitnik/diego/internal/app/api"
	"github.com/ttybitnik/diego/internal/app/core"
	"github.com/ttybitnik/diego/internal/app/domain"
	"github.com/ttybitnik/diego/internal/ports"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI Adapter implements the APIPort and WriterPort interfaces, and cobraModel.
type Adapter struct {
	api   ports.APIPort
	fs    ports.WriterPort
	cobra cobraModel
}

// NewAdapter creates a new instance of the Adapter struct.
func NewAdapter(api ports.APIPort, fs ports.WriterPort) *Adapter {
	return &Adapter{api: api, fs: fs}
}

// cobraAdapter initializes the CLI Adapter.
func cobraAdapter() *Adapter {
	c := core.New()
	fs := filesystem.NewAdapter()
	apia := api.NewApplication(c)
	clia := NewAdapter(apia, fs)

	return clia
}

// cobraModel represents the command-line options for the Diego CLI.
type cobraModel struct {
	format  string
	hugodir string
	input   string
	model   string

	all       bool
	overwrite bool
	scrape    bool
	shortcode bool
}

// global variables pointers for the Cobra/CLI adapter implementation.
var (
	dc  *domain.Core
	dfs *domain.FileSystem
	ca  *Adapter
)

// modelToDomain converts cobraModel into domain.Core and domain.FileSystem.
func modelToDomain(cm cobraModel) {
	dc = &domain.Core{
		All:    cm.all,
		Model:  cm.model,
		Scrape: cm.scrape,
	}

	dfs = &domain.FileSystem{
		Filename:  cm.model,
		Format:    cm.format,
		HugoDir:   cm.hugodir,
		Overwrite: cm.overwrite,
		Shortcode: cm.shortcode,
	}
}

// cobraImport performs the import operation using the API and writes results to filesystem.
func (clia *Adapter) cobraImport(cm cobraModel) {
	err := validateInputs()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	err = checkRequiredDirectories()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	viperSetStringFlags := map[string]*string{
		"format":  &cm.format,
		"hugodir": &cm.hugodir,
	}

	for flag, field := range viperSetStringFlags {
		if viper.IsSet("diego.import." + flag) {
			*field = viper.GetString("diego.import." + flag)
		}
	}

	viperSetBoolFlags := map[string]*bool{
		"all":       &cm.all,
		"overwrite": &cm.overwrite,
		"scrape":    &cm.scrape,
		"shortcode": &cm.shortcode,
	}

	for flag, field := range viperSetBoolFlags {
		if viper.IsSet("diego.import." + flag) {
			*field = viper.GetBool("diego.import." + flag)
		}
	}

	modelToDomain(cm)

	data, err := clia.api.GetImportFile(cm.input, *dc)
	if err != nil {
		log.Fatalln(err)
	}

	err = clia.fs.WriteToFile(data, *dfs)
	if err != nil {
		log.Fatalln(err)
	}

	shortcode, err := clia.api.GetGenerateShortcode(*dc)
	if err != nil {
		log.Fatalln(err)
	}

	err = clia.fs.WriteShortcode(shortcode, *dfs)
	if err != nil {
		log.Fatalln(err)
	}
}

func validateInputs() error {
	hugoDir := "."

	if viper.IsSet("diego.import.hugodir") {
		hugoDir = viper.GetString("diego.import.hugodir")
	}

	hugoDir = filepath.Clean(hugoDir)

	_, err := os.Stat(hugoDir)
	if err != nil {
		return fmt.Errorf("invalid Hugo directory path: %w", err)
	}

	format := domain.OutputYAML

	if viper.IsSet("diego.import.format") {
		format = viper.GetString("diego.import.format")
	}

	formatAllowed := []string{domain.OutputYAML, domain.OutputJSON, domain.OutputTOML, domain.OutputXML}
	formatValid := false
	for _, f := range formatAllowed {
		if format == f {
			formatValid = true
		}
	}

	if !formatValid {
		return fmt.Errorf("invalid output format: %s\nAllowed formats: %v", format, formatAllowed)
	}

	return nil
}

func checkRequiredDirectories() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if viper.IsSet("diego.import.hugodir") {
		currentDir = viper.GetString("diego.import.hugodir")
	}

	currentDir = filepath.Clean(currentDir)

	requiredDirs := []string{
		"archetypes",
		"assets",
		"content",
		"data",
		"layouts",
		"static",
		"themes",
	}

	missingDirs := []string{}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(currentDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			missingDirs = append(missingDirs, dir)
		}
	}

	if len(missingDirs) > 0 {
		return fmt.Errorf(
			"error: invalid Hugo directory structure in: %s\n"+
				"missing the following directories: %v",
			currentDir, missingDirs,
		)
	}

	return nil
}

func bindViperFlags(cmd *cobra.Command, args []string) {
	if !viper.IsSet("diego.import.input") {
		err := cmd.MarkFlagRequired("input")
		if err != nil {
			err = fmt.Errorf("error marking flag as required: %w", err)
			log.Fatalln(err)
		}
	}

	viperBindFlags := []string{"all", "overwrite", "scrape", "shortcode", "format", "hugodir"}

	for _, flag := range viperBindFlags {
		err := viper.BindPFlag("diego.import."+flag, cmd.Flags().Lookup(flag))
		if err != nil {
			err = fmt.Errorf("error binding viper flag: %w", err)
			log.Fatalln(err)
		}
	}
}

func bindInputFlag(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		log.Fatalln("Error binding input flag:", err)
	}
	ca.cobra.input = input

}

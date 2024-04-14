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

// Package filesystem implements the filesystem adapter logic.
package filesystem

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ttybitnik/diego/internal/app/domain"
	"github.com/ttybitnik/diego/internal/app/social"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// Filesystem Adapter implements the FileSystemPort interface.
type Adapter struct{}

// NewAdapter creates a new instance of the Adapter struct.
func NewAdapter() *Adapter {
	return &Adapter{}
}

func (fsa *Adapter) validateWriteToFile(format string) error {
	supportedFormats := map[string]bool{
		domain.OutputYAML: true,
		domain.OutputJSON: true,
		domain.OutputTOML: true,
		domain.OutputXML:  true,
	}

	if !supportedFormats[format] {
		return fmt.Errorf("Error: file format %s not supported", format)
	}

	return nil
}
func (fsa *Adapter) isExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// WriteToFile writes data into filesystem using the given domain.FileSystem configuration.
func (fsa *Adapter) WriteToFile(data []social.Service, dfs domain.FileSystem) error {
	err := fsa.validateWriteToFile(dfs.Format)
	if err != nil {
		log.Fatalln(err)
	}

	var dto []byte

	switch dfs.Format {
	case domain.OutputJSON:
		dto, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Fatalln("Error marshaling into json:", err)
		}
	case domain.OutputTOML:
		tomlData := map[string]interface{}{"Itens": data}
		dto, err = toml.Marshal(tomlData)
		if err != nil {
			log.Fatalln("Error marshaling into toml:", err)
		}
	case domain.OutputXML:
		dto, err = xml.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalln("Error marshaling into xml:", err)
		}
		dto = append([]byte("<Itens>"), dto...)
		dto = append(dto, []byte("</Itens>")...)
	default:
		dto, err = yaml.Marshal(data)
		if err != nil {
			log.Fatalln("Error marshaling into yaml:", err)
		}
	}

	filePath := "data/" + dfs.Filename + "." + dfs.Format
	fileExists := fsa.isExist(filePath)

	if fileExists && !dfs.Overwrite {
		counter := 1
		for {
			newFilePath := fmt.Sprintf("data/%s(%d).%s", dfs.Filename, counter, dfs.Format)
			if !fsa.isExist(newFilePath) {
				filePath = newFilePath
				break
			}
			counter++
		}
	}

	if dfs.HugoDir != "." {
		filePath = dfs.HugoDir + "/" + filePath
	}

	filePath = filepath.Clean(filePath)

	err = os.WriteFile(filePath, dto, 0644)
	if err != nil {
		log.Fatalln("Error writing file into file system:", err)
	}

	fmt.Println("Hugo data file created:", filePath)

	return err
}

// WriteShortcode writes the shortcode into filesystem in the Hugo directory.
func (fsa *Adapter) WriteShortcode(shortcode *string, dfs domain.FileSystem) error {
	if dfs.Shortcode {

		filePath := "layouts/shortcodes/" + dfs.Filename + ".html"

		if dfs.HugoDir != "." {
			filePath = dfs.HugoDir + "/" + filePath
		}

		filePath = filepath.Clean(filePath)

		dirPath := filepath.Dir(filePath)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				log.Fatalln("Error creating [shortcodes] directory:", err)
				return err
			}
		}

		err := os.WriteFile(filePath, []byte(*shortcode), 0644)
		if err != nil {
			log.Fatalln("Error writing shortcode into file system:", err)
		}

		fmt.Println("Hugo shortcode template created:", filePath)

	}

	return nil
}

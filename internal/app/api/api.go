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

// Package api provides a versstile application programming interface.
package api

import (
	"github.com/ttybitnik/diego/internal/app/domain"
	"github.com/ttybitnik/diego/internal/app/social"
)

// Application implements the API Port interface.
type Application struct {
	core Core
}

// NewApplication creates a new instance of Application with the provided core.
func NewApplication(core Core) *Application {
	return &Application{core: core}
}

// GetImportFile delegates the import file operation to the underlying core.
func (a Application) GetImportFile(f string, dc domain.Core) ([]social.Service, error) {
	return a.core.ImportFile(f, dc)
}

// GetGenerateShortcode delegates the shortcode generation to the underlying core.
func (a Application) GetGenerateShortcode(dc domain.Core) (*string, error) {
	return a.core.GenerateShortcode(dc)
}

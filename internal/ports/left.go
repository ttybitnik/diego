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

package ports

import (
	"github.com/ttybitnik/diego/internal/app/domain"
	"github.com/ttybitnik/diego/internal/app/social"
)

// APIPort is the technology neutral port for driver adapters
type APIPort interface {
	GetImportFile(f string, dc domain.Core) ([]social.Service, error)
	GetGenerateShortcode(dc domain.Core) (*string, error)
}

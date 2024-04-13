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

// Package instapaper implements the social interface for instapaper service.
package instapaper

import (
	"fmt"
)

// Default
type List struct {
	Name      string
	URL       string
	Timestamp string
}

func (i *List) BindFile(record *[]string) error {
	i.URL = (*record)[0]
	i.Name = (*record)[1]
	i.Timestamp = (*record)[4]

	return nil
}

func (i *List) FetchFromHTTP() error {
	return nil
}

func (i *List) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// Complete
type ListComplete struct {
	Name      string
	URL       string
	Selection string
	Directory string
	Timestamp string
}

func (i *ListComplete) BindFile(record *[]string) error {
	i.URL = (*record)[0]
	i.Name = (*record)[1]
	i.Selection = (*record)[2]
	i.Directory = (*record)[3]
	i.Timestamp = (*record)[4]

	return nil
}

func (i *ListComplete) FetchFromHTTP() error {
	return nil
}

func (i *ListComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// Common
func htmlPlaylist(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "name" }}
    <tr>
      <td>
	<strong>{{ .name }}</strong>
      </td>
      <td>
	{{ .timestamp }}
      </td>
      <td>
	<a href="{{ .url }}">Link</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

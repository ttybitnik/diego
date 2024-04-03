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

package spotify

import (
	"fmt"
	"strings"
)

// Default
type Library struct {
	Tracks []struct {
		Artist   string `json:"artist"`
		Album    string `json:"album"`
		Track    string `json:"track"`
		URI      string `json:"uri"`
		TrackURL string `json:"trackurl"`
	} `json:"tracks"`
	Albums []struct {
		Artist   string `json:"artist"`
		Album    string `json:"album"`
		URI      string `json:"uri"`
		AlbumURL string `json:"albumurl"`
	} `json:"albums"`
}

func (s *Library) BindFile(record *[]string) error {
	for i := range s.Tracks {
		trackURL, _ := strings.CutPrefix(s.Tracks[i].URI, "spotify:track:")
		s.Tracks[i].TrackURL = "https://open.spotify.com/track/" + trackURL
	}

	for i := range s.Albums {
		albumURL, _ := strings.CutPrefix(s.Albums[i].URI, "spotify:album:")
		s.Albums[i].AlbumURL = "https://open.spotify.com/album/" + albumURL
	}
	return nil
}

func (s *Library) FetchFromHTTP() error {
	return nil
}

func (s *Library) BindHTML(shortcode, comment *string, model string) error {
	return htmlLibrary(shortcode, comment, model)
}

// Complete
type LibraryComplete struct {
	Tracks []struct {
		Artist   string `json:"artist"`
		Album    string `json:"album"`
		Track    string `json:"track"`
		URI      string `json:"uri"`
		TrackURL string `json:"trackurl"`
	} `json:"tracks"`
	Albums []struct {
		Artist   string `json:"artist"`
		Album    string `json:"album"`
		URI      string `json:"uri"`
		AlbumURL string `json:"albumurl"`
	} `json:"albums"`
}

func (s *LibraryComplete) BindFile(record *[]string) error {
	for i := range s.Tracks {
		trackURL, _ := strings.CutPrefix(s.Tracks[i].URI, "spotify:track:")
		s.Tracks[i].TrackURL = "https://open.spotify.com/track/" + trackURL
	}

	for i := range s.Albums {
		albumURL, _ := strings.CutPrefix(s.Albums[i].URI, "spotify:album:")
		s.Albums[i].AlbumURL = "https://open.spotify.com/album/" + albumURL
	}
	return nil
}

func (s *LibraryComplete) FetchFromHTTP() error {
	return nil
}

func (s *LibraryComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlLibrary(shortcode, comment, model)
}

// Common
func htmlLibrary(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s.tracks "track" }}
    <tr>
      <td>
	<strong>{{ .track }}</strong>
      </td>
      <td>
	{{ .artist }}
      </td>
      <td>
	{{ .album }}
      </td>
      <td>
	<a href="{{ .trackurl }}">Spotify</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>
<table>
  <tbody>
    {{ range sort .Site.Data.%s.albums "album" }}
    <tr>
      <td>
	<strong>{{ .album }}</strong>
      </td>
      <td>
	{{ .artist }}
      </td>
      <td>
	<a href="{{ .albumurl }}">Spotify</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model, model)

	return nil
}

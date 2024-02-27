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

package spotify

import (
	"fmt"
	"strings"
)

// Default
type Playlist struct {
	Playlists []struct {
		Name             string `json:"name"`
		LastModifiedDate string `json:"lastModifiedDate"`
		Items            []struct {
			Track struct {
				TrackName  string `json:"trackName"`
				ArtistName string `json:"artistName"`
				AlbumName  string `json:"albumName"`
				TrackURI   string `json:"trackUri"`
				TrackURL   string `json:"trackUrl"`
			} `json:"track"`
			Episode    interface{} `json:"episode"`
			LocalTrack interface{} `json:"localTrack"`
			AddedDate  string      `json:"addedDate"`
		} `json:"items"`
		Description       string `json:"description"`
		NumberOfFollowers int    `json:"numberOfFollowers"`
	} `json:"playlists"`
}

func (s *Playlist) BindFile(record *[]string) error {
	for i := range s.Playlists {
		for j := range s.Playlists[i].Items {
			trackURI := s.Playlists[i].Items[j].Track.TrackURI
			trackURL, _ := strings.CutPrefix(trackURI, "spotify:track:")
			s.Playlists[i].Items[j].Track.TrackURL = "https://open.spotify.com/track/" + trackURL
		}
	}
	return nil
}

func (s *Playlist) FetchFromHTTP() error {
	return nil
}

func (s *Playlist) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// Complete
type PlaylistComplete struct {
	Playlists []struct {
		Name             string `json:"name"`
		LastModifiedDate string `json:"lastModifiedDate"`
		Items            []struct {
			Track struct {
				TrackName  string `json:"trackName"`
				ArtistName string `json:"artistName"`
				AlbumName  string `json:"albumName"`
				TrackURI   string `json:"trackUri"`
				TrackURL   string `json:"trackUrl"`
			} `json:"track"`
			Episode    interface{} `json:"episode"`
			LocalTrack interface{} `json:"localTrack"`
			AddedDate  string      `json:"addedDate"`
		} `json:"items"`
		Description       string `json:"description"`
		NumberOfFollowers int    `json:"numberOfFollowers"`
	} `json:"playlists"`
}

func (s *PlaylistComplete) BindFile(record *[]string) error {
	for i := range s.Playlists {
		for j := range s.Playlists[i].Items {
			trackURI := s.Playlists[i].Items[j].Track.TrackURI
			trackURL, _ := strings.CutPrefix(trackURI, "spotify:track:")
			s.Playlists[i].Items[j].Track.TrackURL = "https://open.spotify.com/track/" + trackURL
		}
	}
	return nil
}

func (s *PlaylistComplete) FetchFromHTTP() error {
	return nil
}

func (s *PlaylistComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// Common
func htmlPlaylist(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s.playlists "name" }}
    <tr>
      <td>
	<strong>{{ .name }}</strong>
      </td>
      <td>
	{{ .items.track.trackname }}
      </td>
      <td>
	{{ .items.track.artistname }}
      </td>
      <td>
	{{ .items.track.albumname }}
      </td>
      <td>
	<a href="{{ .items.track.trackurl }}">Spotify</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

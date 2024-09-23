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

// Package youtube implements the social interface for youtube service.
package youtube

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// Playlist represents the official input data fields.
type Playlist struct {
	Name      string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	ID        string
	Timestamp string
	URL       string
}

// BindFile binds the CSV values into the Playlist struct.
func (y *Playlist) BindFile(record *[]string) error {
	y.ID = (*record)[0]
	y.Timestamp = (*record)[1]
	y.URL = "https://www.youtube.com/watch?v=" + (*record)[0]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the Playlist struct.
func (y *Playlist) FetchFromHTTP() error {
	var name string

	err := fetchPlaylist(y.URL, &name)
	if err != nil {
		return err
	}

	y.Name = name

	return nil
}

// BindHTML generates the Hugo shortcode for the Playlist struct.
func (y *Playlist) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// PlaylistComplete represents the official input data fields.
type PlaylistComplete struct {
	Name      string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	ID        string
	Timestamp string
	URL       string
}

// BindFile binds the CSV values into the PlaylistComplete struct.
func (y *PlaylistComplete) BindFile(record *[]string) error {
	y.ID = (*record)[0]
	y.Timestamp = (*record)[1]
	y.URL = "https://www.youtube.com/watch?v=" + (*record)[0]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the PlaylistComplete struct.
func (y *PlaylistComplete) FetchFromHTTP() error {
	var name string

	err := fetchPlaylist(y.URL, &name)
	if err != nil {
		return err
	}

	y.Name = name

	return nil
}

// BindHTML generates the Hugo shortcode for the PlaylistComplete struct.
func (y *PlaylistComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlPlaylist(shortcode, comment, model)
}

// Common to both structs.
func htmlPlaylist(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "timestamp" }}
    <tr>
      <td>
	<strong>{{ .timestamp }}</strong>
      </td>
      <td>
	{{ .id }}
      </td>
      <td>
	<a href="{{ .url }}">YouTube</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

func fetchPlaylist(url string, name *string) error {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.SetRequestTimeout(15 * time.Second)
	c.IgnoreRobotsTxt = false
	maxRetries := 3
	retryCount := 1

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1000 * time.Millisecond,
	})
	if err != nil {
		log.Fatalf("Error limiting colector: %s", err)
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Fetching \"%s\"...\n", r.URL.String())
	})

	c.OnHTML("meta[name='title']", func(e *colly.HTMLElement) {
		*name = e.Attr("content")
	})

	c.OnError(func(resp *colly.Response, err error) {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			if retryCount <= maxRetries {
				fmt.Printf("Retry %d/%d due to timeout:\n", retryCount, maxRetries)
				_ = resp.Request.Retry()
				retryCount++
			} else {
				log.Printf("Failed fetching \"%s\": %s\n", url, err)
				return
			}
		} else {
			log.Fatalf("Error fetching \"%s\": %s\n", url, err)
		}
	})

	_ = c.Visit(url)

	return nil
}

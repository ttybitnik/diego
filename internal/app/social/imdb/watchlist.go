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

package imdb

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// Default
type Watchlist struct {
	Name      string
	Directors string
	Year      string
	Created   string
	Position  string
	URL       string
	ImgURL    string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
}

func (i *Watchlist) BindFile(record *[]string) error {
	i.Name = (*record)[5]
	i.Directors = (*record)[14]
	i.Year = (*record)[10]
	i.Created = (*record)[2]
	i.Position = (*record)[0]
	i.URL = (*record)[6]

	return nil
}

func (i *Watchlist) FetchFromHTTP() error {
	var imgURL string

	err := fetchWatchlist(i.URL, &imgURL)
	if err != nil {
		return err
	}

	i.ImgURL = imgURL

	return nil
}

func (i *Watchlist) BindHTML(shortcode, comment *string, model string) error {
	return htmlWatchlist(shortcode, comment, model)
}

// Complete
type WatchlistComplete struct {
	Name        string
	Directors   string
	Year        string
	Created     string
	Position    string
	URL         string
	ImgURL      string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Const       string
	Modified    string
	Description string
	TitleType   string
	IMDbRating  string
	RuntimeMins string
	Genres      string
	NumVotes    string
	ReleaseDate string
	YourRating  string
	DateRated   string
}

func (i *WatchlistComplete) BindFile(record *[]string) error {
	i.Name = (*record)[5]
	i.Directors = (*record)[14]
	i.Year = (*record)[10]
	i.Created = (*record)[2]
	i.Position = (*record)[0]
	i.URL = (*record)[6]
	i.Const = (*record)[1]
	i.Modified = (*record)[3]
	i.Description = (*record)[4]
	i.TitleType = (*record)[7]
	i.IMDbRating = (*record)[8]
	i.RuntimeMins = (*record)[9]
	i.Genres = (*record)[11]
	i.NumVotes = (*record)[12]
	i.ReleaseDate = (*record)[13]
	i.YourRating = (*record)[15]
	i.DateRated = (*record)[16]

	return nil
}

func (i *WatchlistComplete) FetchFromHTTP() error {
	var imgURL string

	err := fetchWatchlist(i.URL, &imgURL)
	if err != nil {
		return err
	}

	i.ImgURL = imgURL

	return nil
}

func (i *WatchlistComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlWatchlist(shortcode, comment, model)
}

// Common
func htmlWatchlist(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "name" }}
    <tr>
      <td>
	<strong>{{ .name }}</strong>
      </td>
      <td>
	{{ .directors }}
      </td>
      <td>
	{{ .year }}
      </td>
      <td>
	{{ .created }}
      </td>
      <td>
	{{ .position }}
      </td>
      <td>
	<a href="{{ .url }}">IMDb</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

func fetchWatchlist(url string, imgURL *string) error {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.SetRequestTimeout(15 * time.Second)
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

	c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
		*imgURL = e.Attr("content")
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

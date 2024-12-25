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

// Package imdb implements the social interface for imdb service.
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

// List represents the official input data fields.
type List struct {
	Name         string
	OriginalName string
	Directors    string
	Year         string
	YourRating   string
	URL          string
	ImgURL       string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
}

// BindFile binds the CSV values into the List struct.
func (i *List) BindFile(record *[]string) error {
	i.Name = (*record)[5]
	i.OriginalName = (*record)[6]
	i.Directors = (*record)[15]
	i.Year = (*record)[11]
	i.YourRating = (*record)[16]
	i.URL = (*record)[7]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the List struct.
func (i *List) FetchFromHTTP() error {
	var imgURL string

	err := fetchList(i.URL, &imgURL)
	if err != nil {
		return err
	}

	i.ImgURL = imgURL

	return nil
}

// BindHTML generates the Hugo shortcode for the List struct.
func (i *List) BindHTML(shortcode, comment *string, model string) error {
	return htmlList(shortcode, comment, model)
}

// ListComplete represents the official input data fields.
type ListComplete struct {
	Name         string
	OriginalName string
	Directors    string
	Year         string
	YourRating   string
	URL          string
	ImgURL       string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Position     string
	Const        string
	Created      string
	Modified     string
	Description  string
	TitleType    string
	IMDbRating   string
	RuntimeMins  string
	Genres       string
	NumVotes     string
	ReleaseDate  string
	DateRated    string
}

// BindFile binds the CSV values into the ListComplete struct.
func (i *ListComplete) BindFile(record *[]string) error {
	i.Name = (*record)[5]
	i.OriginalName = (*record)[6]
	i.Directors = (*record)[15]
	i.Year = (*record)[11]
	i.YourRating = (*record)[16]
	i.URL = (*record)[7]
	i.Position = (*record)[0]
	i.Const = (*record)[1]
	i.Created = (*record)[2]
	i.Modified = (*record)[3]
	i.Description = (*record)[4]
	i.TitleType = (*record)[8]
	i.IMDbRating = (*record)[9]
	i.RuntimeMins = (*record)[10]
	i.Genres = (*record)[12]
	i.NumVotes = (*record)[13]
	i.ReleaseDate = (*record)[14]
	i.DateRated = (*record)[17]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the ListComplete struct.
func (i *ListComplete) FetchFromHTTP() error {
	var imgURL string

	err := fetchList(i.URL, &imgURL)
	if err != nil {
		return err
	}

	i.ImgURL = imgURL

	return nil
}

// BindHTML generates the Hugo shortcode for the ListComplete struct.
func (i *ListComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlList(shortcode, comment, model)
}

// Common to both structs.
func htmlList(shortcode, comment *string, model string) error {
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
	{{ .yourrating }}
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

func fetchList(url string, imgURL *string) error {
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

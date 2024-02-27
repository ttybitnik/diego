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

package letterboxd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// Default
type Films struct {
	Name     string
	Director string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Year     string
	URL      string
	ImgURL   string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Date     string
}

func (l *Films) BindFile(record *[]string) error {
	l.Name = (*record)[1]
	l.Year = (*record)[2]
	l.URL = (*record)[3]
	l.Date = (*record)[0]

	return nil
}

func (l *Films) FetchFromHTTP() error {
	var director, imgURL string

	err := fetchFilms(l.URL, &director, &imgURL)
	if err != nil {
		return err
	}

	l.Director = director
	l.ImgURL = imgURL

	return nil
}

func (l *Films) BindHTML(shortcode, comment *string, model string) error {
	return htmlFilms(shortcode, comment, model)
}

// Complete
type FilmsComplete struct {
	Name     string
	Director string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Year     string
	URL      string
	ImgURL   string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Date     string
}

func (l *FilmsComplete) BindFile(record *[]string) error {
	l.Name = (*record)[1]
	l.Year = (*record)[2]
	l.URL = (*record)[3]
	l.Date = (*record)[0]

	return nil
}

func (l *FilmsComplete) FetchFromHTTP() error {
	var director, imgURL string

	err := fetchFilms(l.URL, &director, &imgURL)
	if err != nil {
		return err
	}

	l.Director = director
	l.ImgURL = imgURL

	return nil
}

func (l *FilmsComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlFilms(shortcode, comment, model)
}

// Common
func htmlFilms(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "name" }}
    <tr>
      <td>
	<strong>{{ .name }}</strong>
      </td>
      <td>
	{{ .year }}
      </td>
      <td>
	{{ .date }}
      </td>
      <td>
	<a href="{{ .url }}">Letterboxd</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

func fetchFilms(url string, director, imgURL *string) error {
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

	c.OnHTML("meta[name='twitter:data1']", func(e *colly.HTMLElement) {
		*director = e.Attr("content")
	})

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		scriptContent := e.Text

		split := strings.SplitAfter(scriptContent, ".jpg")
		if len(split) <= 1 {
			log.Println("Scrape: no image found for", url)
		}
		if len(split) >= 2 {
			trim := strings.SplitAfter(split[0], "image\":\"")
			replace := strings.Replace(trim[1], "0-230-0-345-crop", "0-500-0-750-crop", 1)
			*imgURL = replace
		}
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

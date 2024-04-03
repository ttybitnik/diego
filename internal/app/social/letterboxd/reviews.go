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
type Reviews struct {
	Name     string
	Director string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Year     string
	URL      string
	ImgURL   string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Date     string
	Rating   string
	Review   string
}

func (l *Reviews) BindFile(record *[]string) error {
	l.Name = (*record)[1]
	l.Year = (*record)[2]
	l.URL = (*record)[3]
	l.Date = (*record)[0]
	l.Rating = (*record)[4]
	l.Review = (*record)[6]

	return nil
}

func (l *Reviews) FetchFromHTTP() error {
	var director, imgURL string

	err := fetchReviews(l.URL, &director, &imgURL)
	if err != nil {
		return err
	}

	l.Director = director
	l.ImgURL = imgURL

	return nil
}

func (l *Reviews) BindHTML(shortcode, comment *string, model string) error {
	return htmlReviews(shortcode, comment, model)
}

// Complete
type ReviewsComplete struct {
	Name        string
	Director    string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Year        string
	URL         string
	ImgURL      string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
	Date        string
	Rating      string
	Rewatch     string
	Review      string
	Tags        string
	WatchedDate string
}

func (l *ReviewsComplete) BindFile(record *[]string) error {
	l.Name = (*record)[1]
	l.Year = (*record)[2]
	l.URL = (*record)[3]
	l.Date = (*record)[0]
	l.Rating = (*record)[4]
	l.Rewatch = (*record)[5]
	l.Review = (*record)[6]
	l.Tags = (*record)[7]
	l.WatchedDate = (*record)[8]

	return nil
}

func (l *ReviewsComplete) FetchFromHTTP() error {
	var director, imgURL string

	err := fetchReviews(l.URL, &director, &imgURL)
	if err != nil {
		return err
	}

	l.Director = director
	l.ImgURL = imgURL

	return nil
}

func (l *ReviewsComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlReviews(shortcode, comment, model)
}

// Common
func htmlReviews(shortcode, comment *string, model string) error {
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
	{{ .rating }}
      </td>
      <td>
	{{ .review }}
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

func fetchReviews(url string, director, imgURL *string) error {
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

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		scriptContent := e.Text

		directorSplit := strings.SplitAfter(scriptContent, "director\":[{\"@type\":\"Person\",\"name\":\"")
		if len(directorSplit) <= 1 {
			log.Println("Scrape: no director found for", url)
		}
		if len(directorSplit) >= 2 {
			directorTrim := strings.SplitAfter(directorSplit[1], "\"")
			directorCut, _ := strings.CutSuffix(directorTrim[0], "\"")
			*director = directorCut
		}

		imgURLSplit := strings.SplitAfter(scriptContent, ".jpg")
		if len(imgURLSplit) <= 1 {
			log.Println("Scrape: no image found for", url)
		}
		if len(imgURLSplit) >= 2 {
			imgURLTrim := strings.SplitAfter(imgURLSplit[0], "image\":\"")
			imgURLReplace := strings.Replace(imgURLTrim[1], "0-230-0-345-crop", "0-500-0-750-crop", 1)
			*imgURL = imgURLReplace
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

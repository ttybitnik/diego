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

package goodreads

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
type Library struct {
	Name          string
	Author        string
	MyRating      string
	YearPublished string
	URL           string
	ImgURL        string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
}

func (g *Library) BindFile(record *[]string) error {
	g.Name = (*record)[1]
	g.Author = (*record)[2]
	g.MyRating = (*record)[7]
	g.YearPublished = (*record)[12]
	g.URL = "https://www.goodreads.com/book/show/" + (*record)[0]

	return nil
}

func (g *Library) FetchFromHTTP() error {
	var imgURL string

	err := fetchLibrary(g.URL, &imgURL)
	if err != nil {
		return err
	}

	g.ImgURL = imgURL

	return nil
}

func (g *Library) BindHTML(shortcode, comment *string, model string) error {
	return htmlLibrary(shortcode, comment, model)
}

// Complete
type LibraryComplete struct {
	Name               string
	Author             string
	AdditionalAuthors  string
	ISBN               string
	ISBN13             string
	MyRating           string
	AverageRating      string
	Publisher          string
	Binding            string
	NumberOfPages      string
	YearPublished      string
	Year               string
	DateRead           string
	DateAdded          string
	Bookshelves        string
	BookshelvesWithPos string
	ExclusiveShelf     string
	MyReview           string
	Spoiler            string
	PrivateNotes       string
	ReadCount          string
	OwnedCopies        string
	ID                 string
	URL                string
	ImgURL             string `yaml:",omitempty" json:",omitempty" toml:",omitempty" xml:",omitempty"`
}

func (g *LibraryComplete) BindFile(record *[]string) error {
	g.Name = (*record)[1]
	g.Author = (*record)[2]
	g.AdditionalAuthors = (*record)[4]
	g.ISBN = (*record)[5]
	g.ISBN13 = (*record)[6]
	g.MyRating = (*record)[7]
	g.AverageRating = (*record)[8]
	g.Publisher = (*record)[9]
	g.Binding = (*record)[10]
	g.NumberOfPages = (*record)[11]
	g.YearPublished = (*record)[12]
	g.Year = (*record)[13]
	g.DateRead = (*record)[14]
	g.DateAdded = (*record)[15]
	g.Bookshelves = (*record)[16]
	g.BookshelvesWithPos = (*record)[17]
	g.ExclusiveShelf = (*record)[18]
	g.MyReview = (*record)[19]
	g.Spoiler = (*record)[20]
	g.PrivateNotes = (*record)[21]
	g.ReadCount = (*record)[22]
	g.OwnedCopies = (*record)[23]
	g.ID = (*record)[0]
	g.URL = "https://www.goodreads.com/book/show/" + (*record)[0]

	return nil
}

func (g *LibraryComplete) FetchFromHTTP() error {
	var imgURL string

	err := fetchLibrary(g.URL, &imgURL)
	if err != nil {
		return err
	}

	g.ImgURL = imgURL

	return nil
}

func (g *LibraryComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlLibrary(shortcode, comment, model)
}

// Common
func htmlLibrary(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "name" }}
    <tr>
      <td>
	<strong>{{ .name }}</strong>
      </td>
      <td>
	{{ .author }}
      </td>
      <td>
	{{ .myrating }}
      </td>
      <td>
	{{ .yearpublished }}
      </td>
      <td>
	<a href="{{ .url }}">goodreads</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}

func fetchLibrary(url string, imgURL *string) error {
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

	c.OnHTML("div.BookCover__image img", func(e *colly.HTMLElement) {
		*imgURL = e.Attr("src")
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

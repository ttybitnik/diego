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

package cli

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/ttybitnik/diego/internal/app/domain"

	"gopkg.in/yaml.v3"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

var cases = []struct {
	command   string
	fixture   []string
	returnErr []bool
	model     string
	golden    string
}{
	{"GoodreadsLibrary", []string{"library", "invalid", "empty"}, []bool{false, true, true}, domain.GoodreadsLibrary, "library"},
	{"ImdbList", []string{"list", "invalid", "empty"}, []bool{false, true, true}, domain.ImdbList, "list"},
	{"ImdbRatings", []string{"ratings", "invalid", "empty"}, []bool{false, true, true}, domain.ImdbRatings, "ratings"},
	{"ImdbWatchlist", []string{"watchlist", "invalid", "empty"}, []bool{false, true, true}, domain.ImdbWatchlist, "watchlist"},
	{"InstapaperList", []string{"list", "invalid", "empty"}, []bool{false, true, true}, domain.InstapaperList, "list"},
	{"LetterboxdDiary", []string{"diary", "invalid", "empty"}, []bool{false, true, true}, domain.LetterboxdDiary, "diary"},
	{"LetterboxdFilms", []string{"films", "invalid", "empty"}, []bool{false, true, true}, domain.LetterboxdFilms, "films"},
	{"LetterboxdReviews", []string{"reviews", "invalid", "empty"}, []bool{false, true, true}, domain.LetterboxdReviews, "reviews"},
	{"LetterboxdWatchlist", []string{"watchlist", "invalid", "empty"}, []bool{false, true, true}, domain.LetterboxdWatchlist, "watchlist"},
	{"SpotifyLibrary", []string{"library", "invalid", "empty"}, []bool{false, true, true}, domain.SpotifyLibrary, "library"},
	{"SpotifyPlaylist", []string{"playlist", "invalid", "empty"}, []bool{false, true, true}, domain.SpotifyPlaylist, "playlist"},
	{"YoutubePlaylist", []string{"playlist", "invalid", "empty"}, []bool{false, true, true}, domain.YoutubePlaylist, "playlist"},
	{"YoutubeSubscriptions", []string{"subscriptions", "invalid", "empty"}, []bool{false, true, true}, domain.YoutubeSubscriptions, "subscriptions"},
}

func serviceSubdirAndFormat(t *testing.T, command string) (string, string) {
	t.Helper()

	var subdir, format string

	switch {
	case strings.HasPrefix(command, "Goodreads"):
		subdir = "goodreads/"
		format = ".csv"
	case strings.HasPrefix(command, "Imdb"):
		subdir = "imdb/"
		format = ".csv"
	case strings.HasPrefix(command, "Instapaper"):
		subdir = "instapaper/"
		format = ".csv"
	case strings.HasPrefix(command, "Letterboxd"):
		subdir = "letterboxd/"
		format = ".csv"
	case strings.HasPrefix(command, "Spotify"):
		subdir = "spotify/"
		format = ".json"
	case strings.HasPrefix(command, "Youtube"):
		subdir = "youtube/"
		format = ".csv"
	default:
		t.Fatalf("Error finding subdir for %s", command)
	}

	return subdir, format
}

func fixtureValue(t *testing.T, command string, fixture string) string {
	t.Helper()

	subdir, format := serviceSubdirAndFormat(t, command)
	fixturePath := "testdata/" + subdir + fixture + format

	return fixturePath
}

func goldenValue(t *testing.T, command, golden string, gotByte []byte, update bool) string {
	t.Helper()

	subdir, _ := serviceSubdirAndFormat(t, command)
	goldenPath := "testdata/" + subdir + golden + ".golden"

	if update {
		err := os.WriteFile(goldenPath, gotByte, 0644)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", goldenPath, err)
		}

		return string(gotByte)
	}

	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("Error reading file %s: %s", goldenPath, err)
	}

	return string(content)
}

func TestCommands(t *testing.T) {
	ca := cobraAdapter()

	for _, tc := range cases {
		for j := range tc.fixture {
			dc := domain.Core{
				Model:  tc.model,
				All:    false,
				Scrape: false,
			}

			t.Run(tc.command, func(t *testing.T) {
				dto, err := ca.api.GetImportFile(fixtureValue(t, tc.command, tc.fixture[j]), dc)
				gotErr := err != nil
				if gotErr != tc.returnErr[j] {
					t.Fatalf("Unexpected returning err.\nWant: %v, Got: %v", tc.returnErr[j], gotErr)
				}

				if j == 0 {
					var gotByte []byte
					gotByte, err = yaml.Marshal(dto)
					if err != nil {
						t.Fatalf("Error marshaling into yaml: %s", err)
					}

					want := goldenValue(t, tc.command, tc.golden, gotByte, *update)
					got := string(gotByte)

					if got != want {
						t.Fatalf("Unexpected golden value.\nWant:\n%v\nGot:\n%v", want, got)
					}
				}
			})

			dc.All = true

			t.Run(tc.command+"Complete", func(t *testing.T) {
				dto, err := ca.api.GetImportFile(fixtureValue(t, tc.command, tc.fixture[j]), dc)
				gotErr := err != nil
				if gotErr != tc.returnErr[j] {
					t.Fatalf("Unexpected returning err.\nWant: %v, Got: %v", tc.returnErr[j], gotErr)
				}

				if j == 0 {
					var gotByte []byte
					gotByte, err = yaml.Marshal(dto)
					if err != nil {
						t.Fatalf("Error marshaling into yaml: %s", err)
					}

					want := goldenValue(t, tc.command+"Complete", "complete_"+tc.golden, gotByte, *update)
					got := string(gotByte)

					if got != want {
						t.Fatalf("Unexpected golden value.\nWant:\n%v\nGot:\n%v", want, got)
					}
				}
			})
		}
	}
}

func TestShortcodes(t *testing.T) {
	ca := cobraAdapter()

	for _, tc := range cases {
		dc := domain.Core{
			Model:  tc.model,
			All:    false,
			Scrape: false,
		}

		dfs := domain.FileSystem{
			Overwrite: false,
			Shortcode: false,
		}

		dfs.Shortcode = true

		t.Run(tc.command, func(t *testing.T) {
			gotPtr, err := ca.api.GetGenerateShortcode(dc)
			gotErr := err != nil
			if gotErr != tc.returnErr[0] {
				t.Fatalf("Unexpected returning err.\nWant: %v, Got: %v", tc.returnErr[0], gotErr)
			}

			gotByte := []byte(*gotPtr)

			want := goldenValue(t, tc.command, "shortcode_"+tc.golden, gotByte, *update)
			got := string(gotByte)

			if got != want {
				t.Fatalf("Unexpected golden value.\nWant:\n%v\nGot:\n%v", want, got)
			}

		})

		dc.All = true

		t.Run(tc.command+"Complete", func(t *testing.T) {
			gotPtr, err := ca.api.GetGenerateShortcode(dc)
			gotErr := err != nil
			if gotErr != tc.returnErr[0] {
				t.Fatalf("Unexpected returning err.\nWant: %v, Got: %v", tc.returnErr[0], gotErr)
			}

			gotByte := []byte(*gotPtr)

			want := goldenValue(t, tc.command+"Complete", "shortcode_complete_"+tc.golden, gotByte, *update)
			got := string(gotByte)

			if got != want {
				t.Fatalf("Unexpected golden value.\nWant:\n%v\nGot:\n%v", want, got)
			}

		})
	}
}

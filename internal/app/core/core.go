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

// Package core contains the essential application logic.
package core

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ttybitnik/diego/internal/app/domain"
	"github.com/ttybitnik/diego/internal/app/social"
	"github.com/ttybitnik/diego/internal/app/social/goodreads"
	"github.com/ttybitnik/diego/internal/app/social/imdb"
	"github.com/ttybitnik/diego/internal/app/social/instapaper"
	"github.com/ttybitnik/diego/internal/app/social/letterboxd"
	"github.com/ttybitnik/diego/internal/app/social/spotify"
	"github.com/ttybitnik/diego/internal/app/social/youtube"
)

// Core define consrtatins and maximum lengths for the data structs.
const (
	maxAsyncHTTP = 30

	goodreadsLibraryLen     = 24
	imdbListLen             = 17
	imdbRatingsLen          = 13
	imdbWatchlistLen        = 17
	instapaperListLen       = 5
	letterboxdDiaryLen      = 8
	letterboxdFilmsLen      = 4
	letterboxdReviewsLen    = 9
	letterboxdWatchlistLen  = 4
	spotifyLibraryLen       = 8 // TODO: implement JSON len verification.
	spotifyPlaylistLen      = 1 // TODO: implement JSON len verification.
	youtubePlaylistLen      = 2
	youtubeSubscriptionsLen = 3
)

// App represents the main application struct.
type App struct{}

// New creates a new instance of the App struct.
func New() *App {
	return &App{}
}

func (a *App) selectService(dc domain.Core) (social.Service, int, error) {
	modelMap := map[string]struct {
		service         social.Service
		serviceComplete social.Service
		length          int
	}{
		domain.GoodreadsLibrary: {
			&goodreads.Library{},
			&goodreads.LibraryComplete{},
			goodreadsLibraryLen,
		},
		domain.ImdbList: {
			&imdb.List{},
			&imdb.ListComplete{},
			imdbListLen,
		},
		domain.ImdbRatings: {
			&imdb.Ratings{},
			&imdb.RatingsComplete{},
			imdbRatingsLen,
		},
		domain.ImdbWatchlist: {
			&imdb.Watchlist{},
			&imdb.WatchlistComplete{},
			imdbWatchlistLen,
		},
		domain.InstapaperList: {
			&instapaper.List{},
			&instapaper.ListComplete{},
			instapaperListLen,
		},
		domain.LetterboxdDiary: {
			&letterboxd.Diary{},
			&letterboxd.DiaryComplete{},
			letterboxdDiaryLen,
		},
		domain.LetterboxdFilms: {
			&letterboxd.Films{},
			&letterboxd.FilmsComplete{},
			letterboxdFilmsLen,
		},
		domain.LetterboxdReviews: {
			&letterboxd.Reviews{},
			&letterboxd.ReviewsComplete{},
			letterboxdReviewsLen,
		},
		domain.LetterboxdWatchlist: {
			&letterboxd.Watchlist{},
			&letterboxd.WatchlistComplete{},
			letterboxdWatchlistLen,
		},
		domain.SpotifyLibrary: {
			&spotify.Library{},
			&spotify.LibraryComplete{},
			spotifyLibraryLen,
		},
		domain.SpotifyPlaylist: {
			&spotify.Playlist{},
			&spotify.PlaylistComplete{},
			spotifyPlaylistLen,
		},
		domain.YoutubePlaylist: {
			&youtube.Playlist{},
			&youtube.PlaylistComplete{},
			youtubePlaylistLen,
		},
		domain.YoutubeSubscriptions: {
			&youtube.Subscriptions{},
			&youtube.SubscriptionsComplete{},
			youtubeSubscriptionsLen,
		},
	}

	modelSelected, ok := modelMap[dc.Model]
	if !ok {
		return nil, 0, fmt.Errorf("Model type '%s' not valid.", dc.Model)
	}

	if dc.All {
		modelSelected.service = modelSelected.serviceComplete
	}

	return modelSelected.service, modelSelected.length, nil
}

func (a *App) parseFromCSV(reader *csv.Reader, dc domain.Core) ([]social.Service, error) {
	var wg sync.WaitGroup
	scrapeSemaphore := make(chan struct{}, maxAsyncHTTP)
	errCh := make(chan error, 1)

	header, err := reader.Read()
	if len(header) == 0 {
		return nil, fmt.Errorf("Empty CSV file: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Error skipping the first line: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading the CSV file: %w", err)
	}

	results := make([]social.Service, len(records))

	for i, record := range records {
		wg.Add(1)
		go func(resultIdx *social.Service, record []string) {
			defer wg.Done()

			newEntity, modelLen, err := a.selectService(dc)
			if err != nil {
				errCh <- fmt.Errorf("Error selecting service: %w", err)
				return
			}

			entityLen := len(record)
			if entityLen != modelLen {
				errCh <- fmt.Errorf("Invalid CSV format. Want: %d fields, Got: %d fields", modelLen, entityLen)
				return
			}

			err = newEntity.BindFile(&record)
			if err != nil {
				errCh <- fmt.Errorf("Error binding CSV: %w", err)
				return
			}

			if dc.Scrape {
				scrapeSemaphore <- struct{}{}
				defer func() {
					<-scrapeSemaphore
				}()
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				err := newEntity.FetchFromHTTP()
				if err != nil {
					errCh <- fmt.Errorf("Error fetching HTTP: %w", err)
					return
				}
			}

			*resultIdx = newEntity

		}(&results[i], record)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (a *App) parseFromJSON(recorder *json.Decoder, dc domain.Core, data *[]social.Service) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	scrapeSemaphore := make(chan struct{}, maxAsyncHTTP)
	errCh := make(chan error, 1)
	record := []string{}

	if !recorder.More() {
		return fmt.Errorf("Empty JSON file.")
	}

	for {
		newEntity, _, err := a.selectService(dc)
		if err != nil {
			return fmt.Errorf("Error selecting service: %w", err)
		}

		err = recorder.Decode(&newEntity)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error parsing JSON: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			err = newEntity.BindFile(&record)
			if err != nil {
				errCh <- fmt.Errorf("Error binding JSON: %w", err)
			}

			if dc.Scrape {
				scrapeSemaphore <- struct{}{}
				defer func() {
					<-scrapeSemaphore
				}()
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				err = newEntity.FetchFromHTTP()
				if err != nil {
					errCh <- fmt.Errorf("Error fetching HTTP: %w", err)

				}
			}
		}()

		mu.Lock()
		*data = append(*data, newEntity)
		mu.Unlock()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) processDTO(fo *os.File, ext string, dc domain.Core) ([]social.Service, error) {
	data := make([]social.Service, 0)

	switch ext {
	case ".csv":
		var err error
		buffer := bufio.NewReader(fo)
		reader := csv.NewReader(buffer)

		data, err = a.parseFromCSV(reader, dc)
		if err != nil {
			return nil, err
		}
	case ".json":
		buffer := bufio.NewReader(fo)
		decoder := json.NewDecoder(buffer)

		err := a.parseFromJSON(decoder, dc, &data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (a *App) validateImportFile(f string) (string, error) {
	_, err := os.Stat(f)
	if err != nil {
		return "", err
	}

	fe := filepath.Ext(f)

	switch fe {
	case ".csv":
		return fe, nil
	case ".json":
		return fe, nil
	}

	return "", fmt.Errorf("Error wrong file format %s", fe)
}

// ImportFile imports data from input through the specific domain.Core.
func (a *App) ImportFile(f string, dc domain.Core) ([]social.Service, error) {
	ext, err := a.validateImportFile(f)
	if err != nil {
		return nil, err
	}

	fo, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fo.Close()

	fmt.Println("Importing", dc.Model, "from:", f)

	data, err := a.processDTO(fo, ext, dc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GenerateShortcode generates the Hugo shortcode for the given domain.Core.
func (a *App) GenerateShortcode(dc domain.Core) (*string, error) {
	var shortcode string

	newEntity, _, err := a.selectService(dc)
	if err != nil {
		return nil, err
	}

	comment := "<!-- Basic template. " +
		"Read https://gohugo.io/templates/data-templates/ -->"

	err = newEntity.BindHTML(&shortcode, &comment, dc.Model)
	if err != nil {
		return nil, err
	}

	return &shortcode, nil
}

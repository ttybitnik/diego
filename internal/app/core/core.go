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

const (
	maxAsyncHTTP = 30

	goodreadsLibraryLen = 24

	imdbListLen      = 17
	imdbRatingsLen   = 13
	imdbWatchlistLen = 17

	instapaperListLen = 5

	letterboxdDiaryLen     = 8
	letterboxdFilmsLen     = 4
	letterboxdReviewsLen   = 9
	letterboxdWatchlistLen = 4

	spotifyLibraryLen  = 8 // TODO: implement JSON len verification
	spotifyPlaylistLen = 1 // TODO: implement JSON len verification

	youtubePlaylistLen      = 2
	youtubeSubscriptionsLen = 3
)

// Constructor
type App struct{}

func New() *App {
	return &App{}
}

// Logic
func (a *App) selectServiceComplete(dc *domain.Core, sm *social.Service) error {
	switch dc.Model {
	case domain.GoodreadsLibrary:
		*sm = &goodreads.LibraryComplete{}
	case domain.ImdbList:
		*sm = &imdb.ListComplete{}
	case domain.ImdbRatings:
		*sm = &imdb.RatingsComplete{}
	case domain.ImdbWatchlist:
		*sm = &imdb.WatchlistComplete{}
	case domain.InstapaperList:
		*sm = &instapaper.ListComplete{}
	case domain.LetterboxdDiary:
		*sm = &letterboxd.DiaryComplete{}
	case domain.LetterboxdFilms:
		*sm = &letterboxd.FilmsComplete{}
	case domain.LetterboxdReviews:
		*sm = &letterboxd.ReviewsComplete{}
	case domain.LetterboxdWatchlist:
		*sm = &letterboxd.WatchlistComplete{}
	case domain.SpotifyLibrary:
		*sm = &spotify.LibraryComplete{}
	case domain.SpotifyPlaylist:
		*sm = &spotify.PlaylistComplete{}
	case domain.YoutubePlaylist:
		*sm = &youtube.PlaylistComplete{}
	case domain.YoutubeSubscriptions:
		*sm = &youtube.SubscriptionsComplete{}
	default:
		return fmt.Errorf("Model type '%s' not valid.", dc.Model)
	}

	return nil
}

func (a *App) selectService(dc domain.Core) (social.Service, int, error) {
	var sm social.Service
	var mLen int

	switch dc.Model {
	case domain.GoodreadsLibrary:
		sm = &goodreads.Library{}
		mLen = goodreadsLibraryLen
	case domain.ImdbList:
		sm = &imdb.List{}
		mLen = imdbListLen
	case domain.ImdbRatings:
		sm = &imdb.Ratings{}
		mLen = imdbRatingsLen
	case domain.ImdbWatchlist:
		sm = &imdb.Watchlist{}
		mLen = imdbWatchlistLen
	case domain.InstapaperList:
		sm = &instapaper.List{}
		mLen = instapaperListLen
	case domain.LetterboxdDiary:
		sm = &letterboxd.Diary{}
		mLen = letterboxdDiaryLen
	case domain.LetterboxdFilms:
		sm = &letterboxd.Films{}
		mLen = letterboxdFilmsLen
	case domain.LetterboxdReviews:
		sm = &letterboxd.Reviews{}
		mLen = letterboxdReviewsLen
	case domain.LetterboxdWatchlist:
		sm = &letterboxd.Watchlist{}
		mLen = letterboxdWatchlistLen
	case domain.SpotifyLibrary:
		sm = &spotify.Library{}
		mLen = spotifyLibraryLen
	case domain.SpotifyPlaylist:
		sm = &spotify.Playlist{}
		mLen = spotifyPlaylistLen
	case domain.YoutubePlaylist:
		sm = &youtube.Playlist{}
		mLen = youtubePlaylistLen
	case domain.YoutubeSubscriptions:
		sm = &youtube.Subscriptions{}
		mLen = youtubeSubscriptionsLen
	default:
		return nil, 0, fmt.Errorf("Model type '%s' not valid.", dc.Model)
	}

	if dc.All {
		err := a.selectServiceComplete(&dc, &sm)
		if err != nil {
			return nil, 0, err
		}
	}

	return sm, mLen, nil
}

func (a *App) parseFromCSV(reader *csv.Reader, dc domain.Core) ([]social.Service, error) {
	var wg sync.WaitGroup
	var scrapeSemaphore = make(chan struct{}, maxAsyncHTTP)
	errCh := make(chan error, 1)
	emptyModel := make([]social.Service, 0)

	header, err := reader.Read()
	if len(header) == 0 {
		return emptyModel, fmt.Errorf("Empty CSV file: %w", err)
	}
	if err != nil {
		return emptyModel, fmt.Errorf("Error skipping the first line: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return emptyModel, fmt.Errorf("Error reading the CSV file: %w", err)
	}

	results := make([]social.Service, len(records))

	for i, record := range records {
		wg.Add(1)
		go func(resultDest *social.Service, record []string) {
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

			*resultDest = newEntity

		}(&results[i], record)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return results, err
		}
	}

	return results, nil
}

func (a *App) parseFromJSON(recorder *json.Decoder, dc domain.Core, data *[]social.Service) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var scrapeSemaphore = make(chan struct{}, maxAsyncHTTP)
	errCh := make(chan error, 1)
	emptySlice := []string{}

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

			err = newEntity.BindFile(&emptySlice)
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
			return data, err
		}
	case ".json":
		buffer := bufio.NewReader(fo)
		decoder := json.NewDecoder(buffer)

		err := a.parseFromJSON(decoder, dc, &data)
		if err != nil {
			return data, err
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

func (a *App) ImportFile(f string, dc domain.Core) ([]social.Service, error) {
	ext, err := a.validateImportFile(f)
	if err != nil {
		return []social.Service{}, err
	}

	fo, err := os.Open(f)
	if err != nil {
		return []social.Service{}, err
	}
	defer fo.Close()

	fmt.Println("Importing", dc.Model, "from:", f)

	data, err := a.processDTO(fo, ext, dc)
	if err != nil {
		return []social.Service{}, err
	}

	return data, nil
}

func (a *App) GenerateShortcode(dc domain.Core) (*string, error) {
	var shortcode string

	newEntity, _, err := a.selectService(dc)
	if err != nil {
		return &shortcode, err
	}

	comment := "<!-- Basic template. " +
		"Read https://gohugo.io/templates/data-templates/ -->"

	err = newEntity.BindHTML(&shortcode, &comment, dc.Model)
	if err != nil {
		return &shortcode, err
	}

	return &shortcode, nil
}

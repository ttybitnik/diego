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

package domain

const (
	GoodreadsLibrary = "diego_goodreads_library"

	ImdbList      = "diego_imdb_list"
	ImdbRatings   = "diego_imdb_ratings"
	ImdbWatchlist = "diego_imdb_watchlist"

	InstapaperList = "diego_instapaper_list"

	LetterboxdDiary     = "diego_letterboxd_diary"
	LetterboxdFilms     = "diego_letterboxd_films"
	LetterboxdReviews   = "diego_letterboxd_reviews"
	LetterboxdWatchlist = "diego_letterboxd_watchlist"

	SpotifyLibrary  = "diego_spotify_library"
	SpotifyPlaylist = "diego_spotify_playlist"

	YoutubePlaylist      = "diego_youtube_playlist"
	YoutubeSubscriptions = "diego_youtube_subscriptions"

	OutputYAML = "yaml"
	OutputJSON = "json"
	OutputTOML = "toml"
	OutputXML  = "xml"
)

// Core represents the configuration domain for the APIPort.
type Core struct {
	Model  string
	Scrape bool
	All    bool
}

// FileSystem represents the configuration domain for the WriterPort.
type FileSystem struct {
	Filename  string
	Format    string
	HugoDir   string
	Overwrite bool
	Shortcode bool
}

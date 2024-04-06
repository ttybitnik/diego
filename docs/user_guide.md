# Diego User Guide

> [!TIP]
> Utilize the outline button in the top right corner for easy navigation through the headings.

## Commands

Comprehensive list of current available commands with usage examples.

### `diego completion`

Generate the autocompletion script for diego for the specified shell.
See each sub-command's help for details on how to use the generated script.

### `diego completion bash`

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(diego completion bash)

To load completions for every new session, execute once:

**Linux**

	diego completion bash > /etc/bash_completion.d/diego

**macOS**

	diego completion bash > $(brew --prefix)/etc/bash_completion.d/diego

You will need to start a new shell for this setup to take effect.


```
diego completion bash
```

### `diego completion fish`

Generate the autocompletion script for the fish shell.
To load completions in your current shell session:

	diego completion fish | source

To load completions for every new session, execute once:

	diego completion fish > ~/.config/fish/completions/diego.fish

You will need to start a new shell for this setup to take effect.


```
diego completion fish [flags]
```

### `diego completion powershell`

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	diego completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
diego completion powershell [flags]
```

### `diego completion zsh`

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(diego completion zsh)

To load completions for every new session, execute once:

**Linux**

	diego completion zsh > "${fpath[1]}/_diego"

**macOS**

	diego completion zsh > $(brew --prefix)/share/zsh/site-functions/_diego

You will need to start a new shell for this setup to take effect.


```
diego completion zsh [flags]
```

### `diego help`

Help provides help for any command in the application.

```
diego help [command] [flags]
```

### `diego import`

Import data from various services into Hugo.

### `diego import goodreads`

Import data from Goodreads.

```
diego import goodreads {library} -i file [flags]
```

**Examples**

```
diego import goodreads library -i library.csv
diego i g library -i library.csv --all --scrape --shortcode
```

### `diego import imdb`

Import data from IMDb.

```
diego import imdb {list|ratings|watchlist} -i file [flags]
```

**Examples**

```
diego import imdb list -i list.csv
diego i i ratings -i ratings.csv --all --scrape --shortcode
```

### `diego import instapaper`

Import data from Instapaper.

```
diego import instapaper {list} -i file [flags]
```

**Examples**

```
diego import instapaper list -i list.csv
diego i ip list -i list.csv --all --shortcode
```

### `diego import letterboxd`

Import data from Letterboxd.

```
diego import letterboxd {diary|films|reviews|watchlist} -i file [flags]
```

**Examples**

```
diego import letterboxd films -i films.csv
diego i l diary -i diary.csv --all --scrape --shortcode
```

### `diego import spotify`

Import data from Spotify.

```
diego import spotify {library|playlist} -i file [flags]
```

**Examples**

```
diego import spotify library -i library.json
diego i s playlist -i playlist.json --all --shortcode
```

### `diego import youtube`

Import data from YouTube.

```
diego import youtube {playlist|subscriptions} -i file [flags]
```

**Examples**

```
diego import youtube subscriptions -i subscriptions.csv
diego i y playlist -i playlist.csv --all --scrape --shortcode
```

### `diego set`

Set a configuration option.

### `diego set all`

Enable or disable the all flag by default.

```
diego set all {true|false|1|0|enabled|disabled} [flags]
```

**Examples**

```
diego set all true
diego s all false
```

### `diego set defaults`

Restore Diego default settings.

```
diego set defaults [flags]
```

**Examples**

```
diego set defaults true
diego s defaults 1
```

### `diego set format`

Set output format for the Hugo data file (default "yaml").

```
diego set format {yaml|json|toml|xml} [flags]
```

**Examples**

```
diego set format yaml
diego s format json
```

### `diego set hugodir`

Set path to the Hugo directory (default ".").

```
diego set hugodir path [flags]
```

**Examples**

```
diego set hugodir ~/Projects/HugoBlog
diego s hugodir ../Projects/HugoSite
```

### `diego set overwrite`

Enable or disable the overwrite flag by default.

```
diego set overwrite {true|false|1|0|enabled|disabled} [flags]
```

**Examples**

```
diego set overwrite true
diego s overwrite false
```

### `diego set scrape`

Enable or disable the scrape flag by default.

```
diego set scrape {true|false|1|0|enabled|disabled} [flags]
```

**Examples**

```
diego set scrape true
diego s scrape false
```

### `diego set shortcode`

Enable or disable the shortcode flag by default.

```
diego set shortcode {true|false|1|0|enabled|disabled} [flags]
```

**Examples**

```
diego set shortcode true
diego s shortcode false
```

### `diego settings`

Show current settings.

```
diego settings [flags]
```

## Flags

Comprehensive list of current available flags.

### General

#### --config string

Config file (default is $HOME/.config/diego/config.yaml).

#### -h, --help

Help for the command.

#### -v, --version

Version for diego.

### Import

#### --all

Import every available field from CSV/JSON file.

#### --format string

Output format for the Hugo data file (default "yaml").

#### --hugodir string

Path to the Hugo directory (default ".").

#### -i, --input string

Path to the CSV/JSON file (required).

#### --overwrite

Overwrite existent output data file.

#### --scrape

Fetch additional data from CSV/JSON links using HTTP.

> [!IMPORTANT]
> This flag is designed to be as lightweight and minimal as possible. It incorporates hard-coded delays and restricted concurrency, focusing on fetching only missing textual information from the URLs in the input file. Despite these precautions, please be advised not to misuse or spam this flag.

#### --shortcode

Generate a shortcode template for Hugo.

## Configuration
Diego provides persistent configuration via the `diego set` and `diego settings` commands, eliminating the need for manual edits to the configuration file except for backup or version control purposes. To revert to default settings at any time, simply execute `diego set defaults`.

Default configuration, stored in **$HOME/.config/diego/config.yaml**
```yaml
diego:
    import:
        all: false
        format: yaml
        hugodir: .
        overwrite: false
        scrape: false
        shortcode: false
```

## Services and Files

List of current supported services and files with data samples.

### Goodreads

Official Goodreads supported files.

**goodreads_library.csv**
```csv
Id,Name,Author,Author l-f,Additional Authors,ISBN,ISBN13,My Rating,Average Rating,Publisher,Binding,Number of Pages,Year Published,Year,Date Read,Date Added,Bookshelves,Bookshelves with positions,Exclusive Shelf,My Review,Spoiler,Private Notes,Read Count,Owned Copies
517435,Memórias do Subsolo,Fyodor Dostoevsky,"Dostoevsky, Fyodor","Boris Schnaiderman, Fiódor Dostoyevski","=""8573261854""","=""9788573261851""",0,4.18,Editora 34,Paperback,152,2009,1864,,2018/01/03,favorites,favorites (#20),read,,,,1,0
13614950,Ein Brasilianer in Berlin,João Ubaldo Ribeiro,"Ribeiro, João Ubaldo",,"=""3939455040""","=""9783939455042""",0,3.92,Suhrkamp,Paperback,,2017,1994,2020/10/13,2020/10/13,favorites,favorites (#11),read,,,,1,0
25587882,Mastering Emacs,Mickey Petersen,"Petersen, Mickey",,"=""""","=""""",0,4.16,,ebook,314,2022,2015,2023/01/22,2023/01/22,,,read,,,,1,0
```

Optional `--scrape` fields: *ImgUrl*

### IMDb

Official IMDb supported files.

**imdb_list.csv**
```csv
Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors,Your Rating,Date Rated
1,tt6532954,2018-01-16,2018-01-16,,No Intenso Agora,https://www.imdb.com/title/tt6532954/,movie,7.3,127,2017,"Documentary, History",606,2017-02-11,João Moreira Salles,10,2018-05-25
2,tt3262342,2018-01-20,2018-01-20,,Loving Vincent,https://www.imdb.com/title/tt3262342/,movie,7.8,94,2017,"Animation, Drama, Mystery",62477,2017-06-12,"DK Welchman, Hugh Welchman",10,2018-01-20
3,tt2763304,2018-01-16,2018-01-16,,T2 Trainspotting,https://www.imdb.com/title/tt2763304/,movie,7.2,117,2017,"Comedy, Crime, Drama",131395,2017-01-22,Danny Boyle,9,2017-05-28
```

Optional `--scrape` fields: *ImgUrl*

**imdb_ratings.csv**
```csv
Const,Your Rating,Date Rated,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors
tt0287467,10,2016-07-24,Hable con ella,https://www.imdb.com/title/tt0287467/,movie,7.9,112,2002,"Drama, Mystery, Romance",116629,2002-03-15,Pedro Almodóvar
tt0058003,10,2018-08-21,Il deserto rosso,https://www.imdb.com/title/tt0058003/,movie,7.5,117,1964,Drama,17257,1964-09-07,Michelangelo Antonioni
tt0118694,10,2021-05-16,Fa yeung nin wah,https://www.imdb.com/title/tt0118694/,movie,8.1,98,2000,"Drama, Romance",163252,2000-05-20,Kar-Wai Wong
```

Optional `--scrape` fields: *ImgUrl*

**imdb_watchlist.csv**
```csv
Position,Const,Created,Modified,Description,Title,URL,Title Type,IMDb Rating,Runtime (mins),Year,Genres,Num Votes,Release Date,Directors,Your Rating,Date Rated
222,tt5537002,2023-12-25,2023-12-25,,Killers of the Flower Moon,https://www.imdb.com/title/tt5537002/,movie,7.8,206,2023,"Crime, Drama, History, Mystery, Romance, Western",139774,2023-05-20,Martin Scorsese,,
223,tt21027780,2023-12-25,2023-12-25,,Kuolleet lehdet,https://www.imdb.com/title/tt21027780/,movie,7.6,81,2023,"Comedy, Drama",6786,2023-05-23,Aki Kaurismäki,,
224,tt26255088,2023-12-25,2023-12-25,,"Elis & Tom, só tinha de ser com você",https://www.imdb.com/title/tt26255088/,movie,7.9,100,2023,Documentary,240,2023-08-24,"Roberto de Oliveira, Jom Tob Azulay",,
```

Optional `--scrape` fields: *ImgUrl*

### Instapaper

Official Instapaper supported files.

**instapaper_list.csv**
```csv
URL,Title,Selection,Folder,Timestamp
https://www.newyorker.com/magazine/2017/09/04/fernando-pessoas-disappearing-act,Fernando Pessoa’s Disappearing Act,,Starred,1679174385
https://dougseven.com/2014/04/17/knightmare-a-devops-cautionary-tale/,Knightmare: A DevOps Cautionary Tale,,Starred,1694613937
https://www.newyorker.com/magazine/2016/08/08/lauren-collins-learns-to-love-in-french,Love in Translation,,Starred,1679174483
```

### Letterboxd

Official Letterboxd supported files.

**letterboxd_diary.csv**
```csv
Date,Name,Year,Letterboxd URI,Rating,Rewatch,Tags,Watched Date
2023-07-16,My Small Land,2022,https://boxd.it/4wyHj9,5,,,2023-07-15
2023-12-03,Mars One,2022,https://boxd.it/5geIaV,4.5,,,2023-12-03
2023-12-24,As Tears Go By,1988,https://boxd.it/5nUHjJ,5,,,2023-12-24
```

Optional `--scrape` fields: *Director*, *ImgUrl*

**letterboxd_likes_films.csv**
```csv
Date,Name,Year,Letterboxd URI
2021-12-14,Sans Soleil,1983,https://boxd.it/28B8
2021-12-14,Dekalog,1989,https://boxd.it/bv4u
2021-12-14,Blade Runner,1982,https://boxd.it/2bcA
```

Optional `--scrape` fields: *Director*, *ImgUrl*

**letterboxd_reviews.csv**
```csv
Date,Name,Year,Letterboxd URI,Rating,Rewatch,Review,Tags,Watched Date
2024-01-07,Fallen Leaves,2023,https://boxd.it/5wluaZ,4,,"Lorem Ipsum!",,2024-01-06
2024-01-21,"Goodbye, Dragon Inn",2003,https://boxd.it/5E4WzB,4,,"Lorem Ipsum!",,2024-01-20
2024-01-24,Pictures of Ghosts,2023,https://boxd.it/5FB5dd,5,,"Lorem Ipsum!",,2024-01-23
```

Optional `--scrape` fields: *Director*, *ImgUrl*

**letterboxd_watchlist.csv**
```csv
Date,Name,Year,Letterboxd URI
2023-12-14,Vivre Sa Vie,1962,https://boxd.it/28s6
2023-12-14,Mon Oncle,1958,https://boxd.it/2apy
2023-12-14,Days,2020,https://boxd.it/oEE8
```

Optional `--scrape` fields: *Director*, *ImgUrl*

### Spotify

Official Spotify supported files.

**spotify_playlist.json**
```json
{
    "playlists": [
	{
	    "name": "A Sétima",
	    "lastModifiedDate": "2020-02-21",
	    "items": [
		{
		    "track": {
			"trackName": "Le tourbillon",
			"artistName": "Jeanne Moreau",
			"albumName": "Le tourbillon de ma vie (Best Of 2017)",
			"trackUri": "spotify:track:0ZRnvyA5MxbgTCfWRw3YU4"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2018-09-09"
		},
		{
		    "track": {
			"trackName": "Cucurrucucu Paloma - Ao Vivo",
			"artistName": "Caetano Veloso",
			"albumName": "Fina Estampa Ao Vivo",
			"trackUri": "spotify:track:2jIyIMXAHeHSz3ip6ccMFl"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2018-09-10"
		},
		{
		    "track": {
			"trackName": "Preciso Me Encontrar",
			"artistName": "Cartola",
			"albumName": "Raizes Do Samba",
			"trackUri": "spotify:track:1op7nM2R2M6FAU6dSCTRWV"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2018-09-09"
		}
	    ],
	    "description": "O cinema e a música: eternodevir.com&#x2F;a-setima&#x2F;",
	    "numberOfFollowers": 4
	},
	{
	    "name": "Rendez-vous 1408",
	    "lastModifiedDate": "2020-02-21",
	    "items": [
		{
		    "track": {
			"trackName": "Love Will Tear Us Apart - 2010 Remaster",
			"artistName": "Joy Division",
			"albumName": "Substance",
			"trackUri": "spotify:track:1r8oPEXqnhUVgkUkJNqEuF"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2016-11-19"
		},
		{
		    "track": {
			"trackName": "Tear You Apart",
			"artistName": "She Wants Revenge",
			"albumName": "She Wants Revenge",
			"trackUri": "spotify:track:5uahjiKiYXxGlF4GdnKPNv"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2016-11-19"
		},
		{
		    "track": {
			"trackName": "Who Can It Be Now?",
			"artistName": "Men At Work",
			"albumName": "The Essential Men At Work",
			"trackUri": "spotify:track:29r3fDexnrto7WABfpblNH"
		    },
		    "episode": null,
		    "localTrack": null,
		    "addedDate": "2016-11-19"
		}
	    ],
	    "description": "Festa de 5 (-:",
	    "numberOfFollowers": 3
	}
  ]
}
```

**spotify_your_library.json**
```json
{
  "tracks": [
    {
      "artist": "Modena City Ramblers",
      "album": "Tracce Clandestine",
      "track": "Crookedwood Polkas",
      "uri": "spotify:track:4pQF8GvUIT2CKmoJG209ak"
    },
    {
      "artist": "Jorge Ben Jor",
      "album": "A Tabua De Esmeralda",
      "track": "O Homem Da Gravata Florida",
      "uri": "spotify:track:0iPjvpxh6XE9Zo2smrWDhw"
    },
    {
      "artist": "Lou Reed",
      "album": "The Raven",
      "track": "Call on Me (feat. Laurie Anderson)",
      "uri": "spotify:track:2WSnY39N88UvTgt1YhvVFA"
    }
  ],
  "albums": [
    {
      "artist": "SEATBELTS",
      "album": "「COWBOY BEBOP」オリジナルサウンドトラック",
      "uri": "spotify:album:6cYPbwsAFAcddFuGeXMR7l"
    },
    {
      "artist": "Lykke Li",
      "album": "Youth Novels",
      "uri": "spotify:album:65ain97ltDAxldCiOcBtHo"
    },
    {
      "artist": "Zaz",
      "album": "Spotify Sessions",
      "uri": "spotify:album:4Oun04R6kdukRGC8djhzHG"
    }
  ],
  "shows": [
    {
      "name": "AntiCast",
      "publisher": "Half Deaf",
      "uri": "spotify:show:55Aug5UvJC8LSWyzEI7xSR"
    },
    {
      "name": "Decrépitos",
      "publisher": "Decrépitos Podcast",
      "uri": "spotify:show:2epGmU5E47xlTJ4kc5Cevs"
    },
    {
      "name": "ID10T with Chris Hardwick",
      "publisher": "Chris Hardwick",
      "uri": "spotify:show:0Jz4xn4nDRQRjQ5SOyTTuv"
    }
  ],
  "episodes": [
  ],
  "bannedTracks": [
  ],
  "artists": [
  ],
  "bannedArtists": [
  ],
  "other": [
  ]
}
```

### YouTube

Official YouTube supported files.

**youtube_playlist.csv**
```csv
Video ID, Playlist video creation timestamp
eY-eyZuW_Uk,2023-12-07T15:57:57+00:00
28tZ-S1LFok,2023-12-07T18:08:55+00:00
J9qja4KmPQQ,2023-12-07T16:49:34+00:00
```

Optional `--scrape` fields: *Name*

**youtube_subscriptions.csv**
```csv
Channel ID,Channel URL,Channel Title
UC9-y-6csu5WGm29I7JiwpnA,http://www.youtube.com/channel/UC9-y-6csu5WGm29I7JiwpnA,Computerphile
UCy0tKL1T7wFoYcxCe0xjN6Q,http://www.youtube.com/channel/UCy0tKL1T7wFoYcxCe0xjN6Q,Technology Connections
UClcE-kVhqyiHCcjYwcpfj9w,http://www.youtube.com/channel/UClcE-kVhqyiHCcjYwcpfj9w,LiveOverflow
```

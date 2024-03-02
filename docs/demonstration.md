# Demonstration

For illustrative purposes, the following is a basic workflow for importing data from Letterboxd using `diego`.

> [!IMPORTANT]
> For a comprehensive understanding of `diego` functionalities, read the [Diego User Guide](user_guide.md) afterwards.

1. Navigate to the Letterboxd website and download your data archive from `Settings > Data > Export Your Data`.

2. Check which Letterboxd files `diego` currently supports by running `diego import letterboxd -h`.

```txt
$ diego import letterboxd -h
Import data from Letterboxd

Usage:
  diego import letterboxd {diary|films|reviews} -i file [flags]

Aliases:
  letterboxd, l

Examples:
diego import letterboxd films -i films.csv
diego i l diary -i diary.csv --all --scrape --shortcode

Flags:
      --all              import every available field from CSV file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for letterboxd
      --hugodir string   path to the Hugo directory (default ".")
  -i, --input string     path to the CSV file (required)
      --overwrite        overwrite existent output data file
      --scrape           fetch additional data from CSV links using HTTP
      --shortcode        generate a shortcode template for Hugo

Global Flags:
      --config string   config file (default is $HOME/.config/diego/config.yaml)
```

3. As displayed in the usage section, `diego import letterboxd` currently supports three Letterboxd files: `diary`, `films`, and `reviews`. For this demonstration, the `films.csv` containing the record of liked movies is used:

```csv
Date,Name,Year,Letterboxd URI
2021-12-14,Sans Soleil,1983,https://boxd.it/28B8
2021-12-14,Dekalog,1989,https://boxd.it/bv4u
2021-12-14,Blade Runner,1982,https://boxd.it/2bcA
```

4. Run `diego import letterboxd films` passing the path to the `films.csv` in the `-i` (`--input`) flag. By default, `diego` assumes the current working path (`.`) as the Hugo directory. `diego` then generates the Hugo data file under the `data` directory.

```txt
$ diego import letterboxd films -i examples/silo/letterboxd_likes_films.csv
Importing diego_letterboxd_films from: examples/silo/letterboxd_likes_films.csv
Hugo data file created: data/diego_letterboxd_films.yaml
```

```yaml
# data/diego_letterboxd_films.yaml
- name: Sans Soleil
  year: "1983"
  url: https://boxd.it/28B8
  date: "2021-12-14"
- name: Dekalog
  year: "1989"
  url: https://boxd.it/bv4u
  date: "2021-12-14"
- name: Blade Runner
  year: "1982"
  url: https://boxd.it/2bcA
  date: "2021-12-14"
```

4.1 Optionally, to generate an Hugo shortcode template for the data file, append the `--shortcode` flag to the command. This creates the Hugo shortcode template under the `layouts/shortcodes` directory.

```txt
$ diego import letterboxd films -i examples/silo/letterboxd_likes_films.csv --shortcode
Importing diego_letterboxd_films from: examples/silo/letterboxd_likes_films.csv
Hugo data file created: data/diego_letterboxd_films.yaml
Hugo shortcode template created: layouts/shortcodes/diego_letterboxd_films.html
```

```html
<!-- layouts/shortcodes/diego_letterboxd_films.html -->
<!-- Basic template. Read https://gohugo.io/templates/data-templates/ -->
<table>
  <tbody>
    {{ range sort .Site.Data.diego_letterboxd_films "name" }}
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
</table>
```

4.2 Optionally, since the Letterboxd `films.csv` does not include director or poster information, append the `--scrape` flag to the command to fetch this additional data for each entry.

> [!IMPORTANT]
> The `--scrape` flag is designed to be as lightweight and minimal as possible. It only fetches missing textual information from the URLs in the input file. To understand how it works and the specific fields retrieved for each command, read the [Diego User Guide](user_guide.md) before using this flag.

```txt
$ diego import letterboxd films -i examples/silo/letterboxd_likes_films.csv --shortcode --scrape
Importing diego_letterboxd_films from: examples/silo/letterboxd_likes_films.csv
Fetching "https://boxd.it/bv4u"...
Fetching "https://boxd.it/28B8"...
Fetching "https://boxd.it/2bcA"...
Hugo data file created: data/diego_letterboxd_films.yaml
Hugo shortcode template created: layouts/shortcodes/diego_letterboxd_films.html
```

```yaml
# data/diego_letterboxd_films.yaml
- name: Sans Soleil
  director: Chris Marker
  year: "1983"
  url: https://boxd.it/28B8
  imgurl: https://a.ltrbxd.com/resized/sm/upload/f1/gs/in/tp/3L5sW2hGHQmNfVBdOsonSfGzSrN-0-500-0-750-crop.jpg
  date: "2021-12-14"
- name: Dekalog
  director: Krzysztof Kieślowski
  year: "1989"
  url: https://boxd.it/bv4u
  imgurl: https://a.ltrbxd.com/resized/film-poster/2/7/4/1/0/5/274105-the-decalogue-0-500-0-750-crop.jpg
  date: "2021-12-14"
- name: Blade Runner
  director: Ridley Scott
  year: "1982"
  url: https://boxd.it/2bcA
  imgurl: https://a.ltrbxd.com/resized/sm/upload/85/io/38/dz/vfzE3pjE5G7G7kcZWrA3fnbZo7V-0-500-0-750-crop.jpg
  date: "2021-12-14"
```

5. To display the imported data file, call the shortcode in any Hugo page or content:

```
{{< diego_letterboxd_films >}}
```

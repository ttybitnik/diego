## diego import imdb

Import data from IMDb

### Synopsis

Import data from IMDb.

```
diego import imdb {list|ratings|watchlist} -i file [flags]
```

### Examples

```
diego import imdb list -i list.csv
diego i i ratings -i ratings.csv --all --scrape --shortcode
```

### Options

```
      --all              import every available field from CSV file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for imdb
      --hugodir string   path to the Hugo directory (default ".")
  -i, --input string     path to the CSV file (required)
      --overwrite        overwrite existent output data file
      --scrape           fetch additional data from CSV links using HTTP
      --shortcode        generate a shortcode template for Hugo
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/diego/config.yaml)
```

### SEE ALSO

* [diego import](diego_import.md)	 - Import data from various services into Hugo


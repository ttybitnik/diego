## diego import youtube

Import data from YouTube

### Synopsis

Import data from YouTube.

```
diego import youtube {playlist|subscriptions} -i file [flags]
```

### Examples

```
diego import youtube subscriptions -i subscriptions.csv
diego i y playlist -i playlist.csv --all --scrape --shortcode
```

### Options

```
      --all              import every available field from CSV file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for youtube
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


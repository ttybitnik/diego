## diego import letterboxd

Import data from Letterboxd

### Synopsis

Import data from Letterboxd.

```
diego import letterboxd {diary|films|reviews} -i file [flags]
```

### Examples

```
diego import letterboxd films -i films.csv
diego i l diary -i diary.csv --all --scrape --shortcode
```

### Options

```
      --all              import every available field from CSV file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for letterboxd
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


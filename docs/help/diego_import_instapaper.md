## diego import instapaper

Import data from Instapaper

### Synopsis

Import data from Instapaper.

```
diego import instapaper {list} -i file [flags]
```

### Examples

```
diego import instapaper list -i list.csv
diego i ip list -i list.csv --all --shortcode
```

### Options

```
      --all              import every available field from CSV file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for instapaper
      --hugodir string   path to the Hugo directory (default ".")
  -i, --input string     path to the CSV file (required)
      --overwrite        overwrite existent output data file
      --shortcode        generate a shortcode template for Hugo
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/diego/config.yaml)
```

### SEE ALSO

* [diego import](diego_import.md)	 - Import data from various services into Hugo


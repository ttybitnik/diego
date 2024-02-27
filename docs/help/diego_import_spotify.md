## diego import spotify

Import data from Spotify

### Synopsis

Import data from Spotify.

```
diego import spotify {library|playlist} -i file [flags]
```

### Examples

```
diego import spotify library -i library.json
diego i s playlist -i playlist.json --all --shortcode
```

### Options

```
      --all              import every available field from JSON file
      --format string    output format for the Hugo data file (default "yaml")
  -h, --help             help for spotify
      --hugodir string   path to the Hugo directory (default ".")
  -i, --input string     path to the JSON file (required)
      --overwrite        overwrite existent output data file
      --shortcode        generate a shortcode template for Hugo
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/diego/config.yaml)
```

### SEE ALSO

* [diego import](diego_import.md)	 - Import data from various services into Hugo


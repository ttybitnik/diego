## diego completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	diego completion fish | source

To load completions for every new session, execute once:

	diego completion fish > ~/.config/fish/completions/diego.fish

You will need to start a new shell for this setup to take effect.


```
diego completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/diego/config.yaml)
```

### SEE ALSO

* [diego completion](diego_completion.md)	 - Generate the autocompletion script for the specified shell


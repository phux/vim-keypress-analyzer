# Vim Keypress Analyzer

vim-keypress-analyzer parses a vim keypress log file generated by `(n)vim -w <a_log_file>` and aggregates key press counts.

## Features

- [x] count keys pressed by operating mode (NORMAL, INSERT, VISUAL, COMMAND)
- [x] count key frequencies (only when not in INSERT mode)
- [x] identify ALT mappings (e.g. `<m-l>`)
- [x] detect repetitive key sequences as antipatterns (e.g. `jjj` or `dddd`)
  - [ ] detect repetitive multi key sequences like `dwdwdwdw`

### TODO

- [ ] (maybe) `<leader><key>` detection?
- [ ] (maybe maybe) build a vim plugin that logs keys and log on the fly to a structured log format

## Example output

```sh
$ vim-keypress-analyzer --file ~/.nvim_keylog --limit 10 --enable-antipatterns

Vim Keypress Analyzer

Key presses per mode (total: 34335)
│─────────────────│───────│───────────│
│ IDENTIFIER (4)  │ COUNT │ SHARE (%) │
│─────────────────│───────│───────────│
│ insert          │ 16.6K │     48.49 │
│ normal          │ 13.3K │     38.82 │
│ visual          │  3.5K │     10.43 │
│ command         │   779 │      2.27 │
│─────────────────│───────│───────────│

Key presses excluding insert mode (total: 17687)
│──────────────────│───────│───────────│
│ IDENTIFIER (10)  │ COUNT │ SHARE (%) │
│──────────────────│───────│───────────│
│ w                │  2.7K │     15.76 │
│ <space>          │  1.5K │      8.68 │
│ j                │  1.1K │      6.43 │
│ k                │   948 │      5.36 │
│ b                │   903 │      5.11 │
│ d                │   604 │      3.41 │
│ e                │   509 │      2.88 │
│ c                │   463 │      2.62 │
│ i                │   454 │      2.57 │
│ o                │   426 │      2.41 │
│──────────────────│───────│───────────│

Found Antipatterns
│───────────────│───────│───────────────────│─────────────────────────│
│ PATTERN (15)  │ COUNT │ TOTAL KEY PRESSES │ AVG KEYS PER OCCURRENCE │
│───────────────│───────│───────────────────│─────────────────────────│
│ www+          │   234 │              1.2K │ 5.22                    │
│ bbb+          │   112 │               518 │ 4.62                    │
│ ko            │    96 │               192 │ 2.00                    │
│ li            │    63 │               126 │ 2.00                    │
│ kkk+          │    42 │               189 │ 4.50                    │
│ jjj+          │    42 │               181 │ 4.31                    │
│ eee+          │    41 │               220 │ 5.37                    │
│ jO            │    24 │                48 │ 2.00                    │
│ hhh+          │    13 │                77 │ 5.92                    │
│ lll+          │    12 │                70 │ 5.83                    │
│ dddd+         │    11 │                46 │ 4.18                    │
│ xxx+          │     9 │                44 │ 4.89                    │
│ ha            │     4 │                 8 │ 2.00                    │
│ XXX+          │     1 │                 6 │ 6.00                    │
│ BBB+          │     1 │                 3 │ 3.00                    │
│───────────────│───────│───────────────────│─────────────────────────│
```

## Install

### Binary from GitHub release

1. Download the archive for your OS (OSX and Linux, OSX is not tested yet) from the [releases page](https://github.com/phux/vim-keypress-analyzer/releases)
1. Extract the binary `vim-keypress-analyzer` to a directory in your `$PATH`

## Collecting keypresses in vim/nvim

Execute (n)vim with the `-w path/to/logfile` (see `:h -w`) flag to generate
a keypress log file.
Note: the file is only written on exiting (n)vim.

Helpful alias to always log your keys:

```sh
alias n='nvim -w ~/.nvim_keylog "$@"'
# or
alias v='vim -w ~/.vim_keylog "$@"'
```

If you want to split the logs per day, to track progress for example:

```sh
mkdir ~/.vim_logs

alias n='nvim -w ~/.vim_logs/$(date -Idate).log "$@"'
# or
alias v='vim -w ~/.vim_logs/$(date -Idate).log "$@"'
```

## Usage

```sh
$ vim-keypress-analyzer -f <a_log_file>
```

### Optional flags

| Flag                          | Description                                                          | Possible values              | Default         |
|-------------------------------|----------------------------------------------------------------------|------------------------------|-----------------|
| `-l`, `--limit`               | limit the number of key presses displayed                            | any positive int             | `0` (unlimited) |
| `-a`, `--enable-antipatterns` | boolean flag, enable a rudimentary antipattern analysis              | flag is present or not       | false           |
| `-e`, `--exclude-modes`       | comma separated list of modes to be excluded from the key press list | insert,normal,command,visual | `insert`        |

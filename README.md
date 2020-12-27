# Vim Keypress Analyzer

vim-keypress-analyzer parses a vim keypress log file generated by `(n)vim -w <a_log_file>` and aggregates key press counts.

## Features

- [x] count keys pressed by operating mode (NORMAL, INSERT, VISUAL, COMMAND)
- [x] count key frequencies (only when not in INSERT mode)
- [x] identify ALT mappings (e.g. `<m-l>`)
- [ ] key sequence analysis (e.g. `jjj`)
  - [x] naive approach done
  - [x] analyze for antipatterns like `li` (= `a`) or `jjjjjj` (= `5j`)?

### TODO

- [ ] (maybe) `<leader><key>` detection?
- [ ] (maybe maybe) build a vim plugin that logs keys and log on the fly to a structured log format

## Example output

```sh
$ vim-keypress-analyzer -l 10 -enable-antipatterns -f ~/.nvim_keylog

Vim Keypress Analyzer

Key presses per mode (total: 22511)
│─────────────────│───────│───────────│
│ IDENTIFIER (4)  │ COUNT │ SHARE (%) │
│─────────────────│───────│───────────│
│ insert          │ 10.9K │     48.52 │
│ normal          │  8.5K │     37.81 │
│ visual          │  2.5K │     11.44 │
│ command         │   502 │      2.23 │
│─────────────────│───────│───────────│

Key presses in non-INSERT modes (total: 11588)
│──────────────────│───────│───────────│
│ IDENTIFIER (10)  │ COUNT │ SHARE (%) │
│──────────────────│───────│───────────│
│ w                │  1.7K │     15.22 │
│ <space>          │   999 │      8.62 │
│ j                │   815 │      7.03 │
│ k                │   693 │      5.98 │
│ b                │   492 │      4.25 │
│ d                │   411 │      3.55 │
│ e                │   336 │      2.90 │
│ o                │   298 │      2.57 │
│ c                │   279 │      2.41 │
│ i                │   277 │      2.39 │
│──────────────────│───────│───────────│

Antipatterns (naive approach)
│───────────────│───────│
│ PATTERN (14)  │ COUNT │
│───────────────│───────│
│ ww            │   516 │
│ bb            │   218 │
│ kk            │   130 │
│ jj            │   109 │
│ ko            │    83 │
│ li            │    42 │
│ hh            │    32 │
│ jO            │    24 │
│ ee            │    24 │
│ ha            │    18 │
│ ll            │    12 │
│ i<cr>         │     8 │
│ a<cr>         │     1 │
│ WW            │     1 │
│───────────────│───────│
```

## Install

### Binary from GitHub release

1. Download the archive for your OS from the [releases page](https://github.com/phux/vim-keypress-analyzer/releases)
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

Optional:

- Append `-limit <number>` (or short `-l <number`) to limit the number of top keystrokes to show.
- Append `-enable-antipatterns` (or short `-a`) to get a rudimentary antipattern analysis.

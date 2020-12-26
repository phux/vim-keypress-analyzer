# Vim Keypress Analyzer

vim-keypress-analyzer parses a vim keypress log file generated by `(n)vim -w <a_log_file>` and aggregates key press counts.

## Features

- [x] count keys pressed by operating mode (NORMAL, INSERT, VISUAL, COMMAND)
- [x] count key frequencies (only when not in INSERT mode)

### TODO

- [ ] identify ALT mappings (e.g. `<m-l>`)
- [ ] key sequence analysis (e.g. `jjj`)
  - [ ] (maybe) analyze for antipatterns like `li` (= `a`) or `jjjjjj` (= `5j`)?
- [ ] (maybe) `<leader><key>` detection?
- (maybe maybe) build a vim plugin that logs keys and log on the fly to a structured log format

## Example output

```sh
$ vim-keypress-analyzer -l 10 -f ~/.nvim_keylog

Vim Keypress Analyzer

Key presses per mode (total: 1068)
│───────────│───────│───────────│
│ NAME (4)  │ COUNT │ SHARE (%) │
│───────────│───────│───────────│
│ insert    │   684 │     64.04 │
│ normal    │   228 │     21.35 │
│ visual    │   101 │      9.46 │
│ command   │    55 │      5.15 │
│───────────│───────│───────────│

Key presses in non-INSERT modes (total: 384)
│───────────│───────│───────────│
│ KEY (10)  │ COUNT │ SHARE (%) │
│───────────│───────│───────────│
│ w         │    49 │     12.76 │
│ <space>   │    44 │     11.46 │
│ o         │    27 │      7.03 │
│ k         │    18 │      4.69 │
│ <cr>      │    17 │      4.43 │
│ :         │    13 │      3.39 │
│ a         │    13 │      3.39 │
│ <esc>     │    13 │      3.39 │
│ j         │    12 │      3.12 │
│ q         │     9 │      2.34 │
│───────────│───────│───────────│
```

## Install

### Binary from GitHub release

1. Download the archive for your OS from the [releases page](https://github.com/phux/vim-keypress-analyzer/releases)
1. Extract the binary `vim-keypress-analyzer` to a directory in your path

## Collecting keypresses in vim/nvim

Execute (n)vim with the `-w path/to/logfile` flag to generate a keypress log file. Note: the file is only written on exiting (n)vim.

Helpful alias:

```sh
n='nvim -w ~/.nvim_keylog "$@"'
v='vim -w ~/.vim_keylog "$@"'
```

## Usage

```sh
$ vim-keypress-analyzer -f <a_log_file>
```

Optional:

Append `-l <number>` to limit the number of top keystrokes to show.

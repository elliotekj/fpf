# fpf - Fuzzy Prompt Finder

**Search through your Claude Code prompt history across projects.**

fpf provides a fast, interactive TUI for searching and retrieving prompts from your Claude Code conversation history. It lets you fuzzy search through all your past prompts, then copies your selection to the clipboard.

[![asciicast](https://asciinema.org/a/0NGUt8MjNaeBkIJbSfUciUVUU.svg)](https://asciinema.org/a/0NGUt8MjNaeBkIJbSfUciUVUU)

## Features

- **Fuzzy search** - Find prompts even with typos or partial matches
- **Project filtering** - Narrow results to a specific project directory using `%p`
- **Preview mode** - View full multi-line prompts before selecting
- **Clipboard integration** - Selected prompts are automatically copied
- **Smart deduplication** - Keeps only the most recent version of duplicate prompts
- **Time awareness** - Shows how long ago each prompt was used
- **Fast** - Efficiently scans and searches large prompt histories

## Installation

### Pre-built binaries

Download the latest release for your platform from the [releases page](https://github.com/elliotekj/fpf/releases).

Extract the archive and move the binary to somewhere in your `$PATH`:

```bash
tar -xzf fpf_*_*.tar.gz
mv fpf /usr/local/bin/
```

### From source

```bash
git clone https://github.com/elliotekj/fpf.git
cd fpf
just build
mv bin/fpf /usr/local/bin/
```

## License

`fpf` is released under the [`Apache License
2.0`](https://github.com/elliotekj/doubly_linked_list/blob/main/LICENSE).

## About

This tool was written by [Elliot Jackson](https://elliotekj.com).

- Blog: [https://elliotekj.com](https://elliotekj.com)
- Email: elliot@elliotekj.com

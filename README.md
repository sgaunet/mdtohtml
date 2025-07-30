[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/mdtohtml)](https://goreportcard.com/report/github.com/sgaunet/mdtohtml)
[![GitHub release](https://img.shields.io/github/release/sgaunet/mdtohtml.svg)](https://github.com/sgaunet/mdtohtml/releases/latest)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/mdtohtml/coverage-badge.svg)
[![coverage](https://github.com/sgaunet/mdtohtml/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/mdtohtml/actions/workflows/coverage.yml)
[![Snapshot Build](https://github.com/sgaunet/mdtohtml/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/mdtohtml/actions/workflows/snapshot.yml)
[![Release Build](https://github.com/sgaunet/mdtohtml/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/mdtohtml/actions/workflows/release.yml)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/mdtohtml/total)

# Markdown to HTML cmd-line tool

A powerful command-line tool to convert Markdown files to HTML with GitHub-style CSS. Built with Go and featuring batch processing, validation, and shell completion support.

## Features

- **Single file conversion** - Convert individual Markdown files to HTML
- **Batch processing** - Convert multiple files at once with pattern matching
- **Recursive processing** - Process entire directory trees
- **Validation** - Check Markdown syntax without generating output
- **Shell completion** - Auto-completion for bash and zsh
- **GitHub-style CSS** - Beautiful GitHub-inspired styling
- **Smart typography** - Optional smart quotes, dashes, and fractions

## Quick Start

```bash
# Convert a single file
mdtohtml README.md README.html

# Batch convert all .md files in a directory
mdtohtml batch ./docs --out-dir ./html

# Validate Markdown syntax
mdtohtml validate document.md
```


## Usage

### Convert Command (Default)

Convert a single Markdown file to HTML:

```bash
# Default command (can omit 'convert')
mdtohtml input.md output.html

# Explicit convert command
mdtohtml convert input.md output.html

# With options
mdtohtml input.md output.html --smartypants=false --latexdashes=false
```

**Options:**
- `--smartypants` (default: true) - Apply smart typography substitutions
- `--latexdashes` (default: true) - Use LaTeX-style dash rules
- `--fractions` (default: true) - Convert fractions like 1/2 to ½

### Batch Command

Convert multiple Markdown files at once:

```bash
# Convert all .md files in a directory
mdtohtml batch ./docs --out-dir ./html

# Use custom file pattern
mdtohtml batch ./docs --pattern "*.markdown" --out-dir ./public

# Process directories recursively
mdtohtml batch ./docs --recursive --out-dir ./output

# With typography options
mdtohtml batch ./docs --out-dir ./html --smartypants=false
```

**Options:**
- `-o, --out-dir` (default: ".") - Output directory for HTML files
- `-p, --pattern` (default: "*.md") - File pattern to match
- `-r, --recursive` - Process directories recursively
- Plus all typography options from convert command

### Validate Command

Check Markdown syntax without generating output:

```bash
# Validate a single file
mdtohtml validate document.md

# Validate with specific parser settings
mdtohtml validate document.md --smartypants=false
```

Returns exit code 0 if valid, non-zero if invalid.

### Shell Completion

Enable auto-completion for your shell:

```bash
# Bash
mdtohtml completion bash > /etc/bash_completion.d/mdtohtml

# Zsh
mdtohtml completion zsh > /usr/local/share/zsh/site-functions/_mdtohtml

# Fish
mdtohtml completion fish > ~/.config/fish/completions/mdtohtml.fish

# PowerShell
mdtohtml completion powershell > mdtohtml.ps1
```

## Examples

```bash
# Convert README to HTML
mdtohtml README.md README.html

# Batch convert documentation
mdtohtml batch ./docs --out-dir ./website/docs --recursive

# Validate before committing
mdtohtml validate *.md

# Convert with plain typography
mdtohtml input.md output.html --smartypants=false --fractions=false
```

# Docker Image

There is a docker image to integrate the binary into your own docker image for example.

For example, the Dockerfile should look like :

```dockerfile
FROM sgaunet/mdtohtml:latest AS mdtohtml

FROM <BASE-IMAGE:VERSION>
...
COPY --from=mdtohtml /usr/bin/mdtohtml /usr/bin/mdtohtml
...
```

## Supported Markdown Features

mdtohtml uses the [Goldmark](https://github.com/yuin/goldmark) parser with the following extensions:

- **GitHub Flavored Markdown (GFM)** - Tables, strikethrough, autolinks, task lists
- **Definition Lists** - Support for `<dl>`, `<dt>`, `<dd>` elements
- **Footnotes** - Reference-style footnotes
- **Typographer** - Smart quotes, dashes, fractions (when enabled)
- **Auto heading IDs** - Automatic generation of heading anchors
- **Unsafe HTML** - Raw HTML is preserved

# Install

## With homebrew

```
brew tap sgaunet/homebrew-tools
brew install sgaunet/tools/mdtohtml
```

## Download release

Download the latest release from the [releases page](https://github.com/sgaunet/mdtohtml/releases) and copy it to `/usr/local/bin` or any directory in your PATH.

## Build from source

```bash
# Clone the repository
git clone https://github.com/sgaunet/mdtohtml.git
cd mdtohtml

# Build with task (recommended)
task build

# Or build directly with Go
go build .

# Install to your PATH
sudo cp mdtohtml /usr/local/bin/
```

### Development

```bash
# Run linter
task linter

# Create snapshot build
task snapshot

# Create release
task release
```
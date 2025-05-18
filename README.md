[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/mdtohtml)](https://goreportcard.com/report/github.com/sgaunet/mdtohtml)
[![GitHub release](https://img.shields.io/github/release/sgaunet/mdtohtml.svg)](https://github.com/sgaunet/mdtohtml/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/mdtohtml)](https://goreportcard.com/report/github.com/sgaunet/mdtohtml)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/mdtohtml/coverage-badge.svg)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/mdtohtml/total)

# Markdown to HTML cmd-line tool

Tool to convert markdown file to html with a css like github.

# Forked project

I clean some code, remove some options and add [the github-markdown CSS](https://github.com/sindresorhus/github-markdown-css)

You can use the README ini tst folder to test the app.

```
mdtohtml README.md README.html
```


# Docker Image

There is a docker image to integrate the binary into your own docker image for example.

For example, the Dockerfile should look like :

```
FROM sgaunet/mdtohtml:0.3.1 AS mdtohtml

FROM <BASE-IMAGE:VERSION>
...
COPY --from=mdtohtml /usr/bin/mdtohtml /usr/bin/mdtohtml
...

```

# Install

## With homebrew

```
brew tap sgaunet/homebrew-tools
brew install sgaunet/tools/mdtohtml
```

## Download release

And copy it to /usr/local/bin
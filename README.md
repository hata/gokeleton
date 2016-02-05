# gokeleton

Skeleton Generator

## Description

Copy and replace keywords to generate skeleton files from local or github repository written by Go

## Usage

```bash
gokeleton -p "key1=value1,key2=value2" <src-path> <dest-path>
```

For example, copy from github repository

```bash
gokeleton -p "key=value" https://github.com/hata/gokeleton /tmp/test
```

Copy from a local directory

```bash
gokeleton -p "key=value" /local/template/path /tmp/test
```

Replace 'key' with 'value' if these keys are found in files.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/hata/gokeleton
```

## Contribution

1. Fork ([https://github.com/hata/gokeleton/fork](https://github.com/hata/gokeleton/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[hata](https://github.com/hata)

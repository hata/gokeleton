package main

import (
    "io"
    "strings"
)

type FileSource interface {
    Name() string
    Reader() (io.ReadCloser, error)
}

type FileSourceFunc func(fileSource FileSource) error

type SourceAccess interface {
    EachSource(callback FileSourceFunc) error
}

type StartParams struct {
    Keywords string
    KeySeparator string
    Arguments []string
}

func StartMain(sp StartParams) error {
    if sp.KeySeparator == "" {
        sp.KeySeparator = ","
    }
    return NewFileAccess(sp.Arguments[0], sp.Arguments[1], stringsToMap(sp.Keywords, sp.KeySeparator)).Run()
}

func stringsToMap(keywords string, sep string) (keyMap map[string]string) {
    var key, value string
    keyMap = map[string]string{}

    if len(keywords) == 0 {
        return
    }

    pairs := strings.Split(keywords, sep)

    for _, kv := range pairs {
        index := strings.IndexByte(kv, '=')
        if index >= 0 {
            key = kv[0:index]
            value = kv[index + 1:]
        } else {
            key = kv
            value = ""
        }

        keyMap[key] = value
    }

    return
}


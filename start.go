package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

type ReplaceFunc func(srcSubPath string, srcContents string) (subPath string, contents string, err error)

type FileSource interface {
    SubPath() string
    IsDir() bool
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
    IncludeSuffixes string
    ExcludeSuffixes string
}

func StartMain(sp StartParams) error {
    srcPath := sp.Arguments[0]
    destPath := sp.Arguments[1]
    keyMap := stringsToMap(sp.Keywords, sp.KeySeparator)
    includeSuffixes := toList(sp.IncludeSuffixes, sp.KeySeparator)
    excludeSuffixes := toList(sp.ExcludeSuffixes, sp.KeySeparator)

    _, err := os.Stat(destPath)
    if os.IsNotExist(err) {
        sa := newSourceAccess(srcPath)
        return copyEachFileSource(sa, destPath, includeSuffixes, excludeSuffixes, newReplaceFunc(keyMap))
    } else if err != nil {
        return err
    } else {
        fmt.Fprintln(os.Stderr, "Error: dest path:", destPath, " already exists")
        return os.ErrExist
    }
}

func newSourceAccess(srcPath string) SourceAccess {
    if strings.Index(srcPath, "http://") == 0 || strings.Index(srcPath, "https://") == 0 {
        return newGithubAccess(srcPath)
    } else {
        return newFileAccess(srcPath)
    }
}

func copyEachFileSource(sa SourceAccess, destPath string, includeSuffixes []string, excludeSuffixes []string, handler ReplaceFunc) error {
    return sa.EachSource(func(fileSource FileSource) error {
        var isDestDir bool
        var contentBytes []byte
        var subPath, contents string

        if fileSource.IsDir() {
            subPath, _, err := handler(fileSource.SubPath(), "")
            if err != nil {
                return err
            }
            return os.MkdirAll(normalizePath(destPath, true) + subPath, 0777)
        }

        reader, err := fileSource.Reader()
        if err != nil {
            return err
        }
        defer reader.Close()

        contentBytes, err = ioutil.ReadAll(reader)
        if err != nil {
            return err
        }

        if isMatchSuffixes(includeSuffixes, fileSource.SubPath()) &&
          !isMatchSuffixes(excludeSuffixes, fileSource.SubPath()) {
            subPath, contents, err = handler(fileSource.SubPath(), string(contentBytes))
            if err != nil {
                return err
            }
            contentBytes = nil
        } else {
            subPath = fileSource.SubPath()
        }

        // This is expected to be created before calling here.
        // Or, ignore error for a dest file is used.
        isDestDir, err = isDirectory(destPath)

        newFilePath := normalizePath(destPath, isDestDir) + subPath
        out, outErr := os.Create(newFilePath)
        if outErr != nil {
            return outErr
        }

        if contentBytes != nil {
            _, err = out.Write(contentBytes)
        } else {
            _, err = out.WriteString(contents)
        }

        if err == nil {
            fmt.Println("Create", newFilePath)
        }

        return err
    })
}

func isMatchSuffixes(suffixes []string, name string) bool {
    for _, suffix := range suffixes {
        if suffix == "*" || strings.HasSuffix(name, suffix) {
            return true
        }
    }
    return false
}

func newReplaceFunc(keywords map[string]string) ReplaceFunc {
    return func (srcSubPath string, srcContents string) (subPath string, contents string, err error) {
        subPath = srcSubPath
        contents = srcContents

        for key, val := range keywords {
            subPath = strings.Replace(subPath, key, val, -1)
            contents = strings.Replace(contents, key, val, -1)
        }

        err = nil
        return
    }
}

func toSubPath(basePath string, fullPath string) string {
    return fullPath[len(basePath):]
}

func normalizePath(path string, isDir bool) string {
    path = filepath.Clean(path)
    if isDir && !strings.HasSuffix(path, pathSep()) {
        path = path + pathSep()
    }
    return path
}

func pathSep() string {
    return string([]byte{os.PathSeparator})
}

func isDirectory(path string) (isDir bool, err error) {
    fInfo, err := os.Stat(path)
    if err != nil {
        return false, err
    }
    return fInfo.IsDir(), nil
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

func toList(listString string, sep string) []string {
    list := strings.Split(listString, sep)
    result := make([]string, len(list), len(list))
    for i, v := range list {
        result[i] = strings.TrimSpace(v)
    }
    return result
}

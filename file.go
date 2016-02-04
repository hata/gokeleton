package main

import (
    "fmt"
    "io"
	"os"
    "path/filepath"
)

type fileAccess struct {
	srcPath string
}

type fileSource struct {
    fullPath string
    subPath string
    info os.FileInfo
}

func newFileAccess(srcPath string) (fa *fileAccess) {
	fa = new(fileAccess)
    isDir, _ := isDirectory(srcPath)
	fa.srcPath = normalizePath(srcPath, isDir)
	return
}

func newFileSource(srcPath string, subPath string, info os.FileInfo) (fs *fileSource) {
    fs = new(fileSource)
    fs.fullPath = srcPath
    fs.subPath = subPath
    fs.info = info
    return
}

// SourceAccess
func (fa *fileAccess) EachSource(callback FileSourceFunc) error {
    return filepath.Walk(fa.srcPath, func (fullPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        subPath := toSubPath(fa.srcPath, fullPath)
        err = callback(newFileSource(fullPath, subPath, info))
        if err != nil {
            fmt.Println("EachSource return error:", err)
        }
        return err
    })
}

// FileSource
func (fs *fileSource) SubPath() string {
    return fs.subPath
}

func (fs *fileSource) IsDir() bool {
    return fs.info.IsDir()
}

func (fs *fileSource) Reader() (io.ReadCloser, error) {
    return os.Open(fs.fullPath)
}

package main

import (
    "errors"
    "fmt"
    "io"
    "math"
	"os"
    "path/filepath"
    "strings"
)

// git access API
// https://developer.github.com/v3/repos/contents/

type ReplaceFunc func(srcSubPath string, srcContents string) (subPath string, contents string, err error)

type FileAccess interface {
    Run() error
}

type fileAccess struct {
	srcPath string
    destPath string
    templateSuffixes []string
    paramMap map[string]string
}


func NewFileAccess(srcPath string, destPath string, paramMap map[string]string) FileAccess {
	return newFileAccess(srcPath, destPath, paramMap)
}

func newFileAccess(srcPath string, destPath string, paramMap map[string]string) (fa *fileAccess) {
	fa = new(fileAccess)
    isDir, _ := isDirectory(srcPath)
	fa.srcPath = normalizePath(srcPath, isDir)
    fa.destPath = normalizePath(destPath, isDir)
    fa.templateSuffixes = []string{".tmpl"}
    fa.paramMap = paramMap
	return
}

func (fa *fileAccess) Run() error {
    _, err := os.Stat(fa.destPath)
    if os.IsNotExist(err) {
        return fa.walkFiles(fileAccessHandler)
    } else if err != nil {
        return err
    } else {
        fmt.Fprintln(os.Stderr, "Error: dest path:", fa.destPath, " already exists")
        return os.ErrExist
    }
}

func (fa *fileAccess) newWalkFunc(handler ReplaceFunc) filepath.WalkFunc {
    return func (fullPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        subPath := toSubPath(fa.srcPath, fullPath)

        if info.IsDir() {
            subPath, _, err = handler(toSubPath(fa.srcPath, fullPath), "")
            if err != nil {
                return err
            }
            return createDirectory(fa.destPath, subPath)
        }

        if isMatchSuffixes(fa.templateSuffixes, info.Name()) {
            contents, err2 := readFile(fullPath, info)
            if err2 == nil {
                subPath, contents, err = handler(subPath, contents)
                if err != nil {
                    return err
                }
                out, outErr := os.Create(fa.destPath + subPath)
                if outErr != nil {
                    return outErr
                }
                _, writeErr := out.WriteString(contents)
                return writeErr
            }
        }

        return copyFile(fullPath, fa.destPath + subPath)
    }
}

func (fa *fileAccess) walkFiles(handler ReplaceFunc) (err error) {
    return filepath.Walk(fa.srcPath, fa.newWalkFunc(handler))
}

func fileAccessHandler (srcSubPath string, srcContents string) (subPath string, contents string, err error) {
    subPath = srcSubPath
    contents = srcContents
    err = nil
    return
}

func isDirectory(path string) (isDir bool, err error) {
    fInfo, err := os.Stat(path)
    if err != nil {
        return false, err
    }
    return fInfo.IsDir(), nil
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


func readFile(fullPath string, fInfo os.FileInfo) (string, error) {
    length64 := fInfo.Size()
    var n int

    if math.MaxInt32 < length64 {
        return "", errors.New("File size is too big")
    }
    length := int(length64)
    buf := make([]byte, length)

    f, err := os.Open(fullPath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    n, err = f.Read(buf)
    if err != nil {
        return "", err
    }

    if err != io.EOF && n != length {
        return "", errors.New("Read few length.")
    }

    return string(buf), nil
}

func createDirectory(parentPath string, subPath string) error {
    fmt.Printf("create dir: %s + %s\n", parentPath, subPath)
    return os.MkdirAll(parentPath + subPath, 0777)
}

func copyFile(srcFilePath string, destFilePath string) error {
    fmt.Printf("copy: from: %s to: %s\n", srcFilePath, destFilePath)

    in, inErr := os.Open(srcFilePath)
    if inErr != nil {
        return inErr
    }
    defer in.Close()

    out, outErr := os.Create(destFilePath)
    if outErr != nil {
        return outErr
    }
    defer out.Close()

    _, copyErr := io.Copy(out, in)
    if copyErr != nil {
        return copyErr
    }

    return nil
}

func toSubPath(basePath string, fullPath string) string {
    return fullPath[len(basePath):]
}

func isMatchSuffixes(suffixes []string, name string) bool {
    for _, suffix := range suffixes {
        if strings.HasSuffix(name, suffix) {
            return true
        }
    }
    return false
}

package main

import (
    "os"
    "testing"
)

func Test_newFileAccess(t *testing.T) {
    curDir, _ := os.Getwd()
    fa := newFileAccess(curDir)
    normalizedCurDir := normalizePath(curDir, true)
    if fa.srcPath != normalizedCurDir {
        t.Error("srcPath should be set correctly")
    }
}

func Test_FileAccess_EachSource_Dir(t *testing.T) {
    curDir, _ := os.Getwd()
    fa := newFileAccess(curDir)
    count := 0

    fa.EachSource(func(fs FileSource) error {
        count++
        return nil
    })

    if count < 5 {
        t.Error("file EachSource failed.")
    }
}

func Test_FileAccess_EachSource_File(t *testing.T) {
    curDir, _ := os.Getwd()
    fa := newFileAccess(curDir + "/file_test.go")
    count := 0
    subPath := "x"
    isDir := false

    fa.EachSource(func(fs FileSource) error {
        count++
        subPath = fs.SubPath()
        isDir = fs.IsDir()
        return nil
    })

    if count != 1 {
        t.Error("file EachSource failed.")
    }

    if subPath != "" {
        t.Error("file should not have a sub path")
    }

    if isDir {
        t.Error("file should return false for IsDir")
    }
}

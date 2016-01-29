package main

import (
    "os"
    "path/filepath"
    "testing"
)

func assertString(t *testing.T, desc string, expected string, result string) {
    if expected != result {
        t.Log("expected: %s, result: %s\n", expected, result)
        t.Error(desc)
    }
}


func Test_isDirectory_not_found(t *testing.T) {
    dir, err := isDirectory("/tmp_not_found")
    if dir {
        t.Error("Failed to validate there is a dir or not")
    }
    if err == nil {
        t.Error("Error should be found %s", err)
    }
}

func Test_isDirectory_found_dir(t *testing.T) {
   dir, err := isDirectory("/tmp")
   if !dir {
       t.Error("/tmp should be found.")
   }
   if err != nil {
       t.Error("/tmp should be found.")
   }
}

func Test_isDirectory_found_file(t *testing.T) {
    curDir, _ := os.Getwd()
    curFile := curDir + "/file_test.go"
    dir, err := isDirectory(curFile)
    if dir {
        t.Error("Current file should be exist. %s", curFile)
    }
    if err != nil {
        t.Error("Current file should be found. %s", curFile)
    }
}

func Test_readFile(t *testing.T) {
    curDir, _ := os.Getwd()
    path := curDir + "/file_test.go"
    fInfo, _ := os.Stat(path)
    contents, _ := readFile(path, fInfo)
    if len(contents) != 0 && len(contents) != int(fInfo.Size()) {
        t.Error("Read file size is wrong")
    }
}

func Test_normalizePath_dir(t *testing.T) {
    path := normalizePath("/tmp", true)
    assertString(t, "normalizePath should have / at the end of a directory path", "/tmp/", path)
}

func Test_normalizePath_file(t *testing.T) {
    curDir, _ := os.Getwd()
    path := normalizePath(curDir + "/file_test.go", false)
    path = normalizePath(path, false)
    assertString(t, "normalizePath should not add / at the end of a file path", "file_test.go", filepath.Base(path))
}

func Test_newFileAccess(t *testing.T) {
    curDir, _ := os.Getwd()
    destPath := "/tmp"
    fa := newFileAccess(curDir, destPath, map[string]string{})
    normalizedCurDir := normalizePath(curDir, true)
    normalizedDestPath := normalizePath(destPath, true)
    assertString(t, "srcPath should be set correctly", normalizedCurDir, fa.srcPath)
    assertString(t, "destPath should be set correctly", normalizedDestPath, fa.destPath)
}

func Test_newFileAccess_with_params(t *testing.T) {
    curDir, _ := os.Getwd()
    destPath := "/tmp"
    fa := newFileAccess(curDir, destPath, map[string]string{"foo":"bar"})
    if fa.paramMap["foo"] != "bar" {
        t.Error("Verify paramMap is set correctly")
    }
}

func Test_Run(t *testing.T) {
    curDir, _ := os.Getwd()
    destPath := "/tmp"
    fa := NewFileAccess(curDir, destPath, map[string]string{})
    if fa == nil {
        t.Error("Failed to create a new FileAccess instance")
    }
}

func Test_toSubPath(t *testing.T) {
    path := toSubPath("/tmp/", "/tmp/foo")
    if path != "foo" {
        t.Error("Verify sub path is an expected value")
    }
}

func Test_isMatchSuffixes_matched(t *testing.T) {
    if !isMatchSuffixes([]string{".foo", ".txt"}, "/tmp/test.txt") {
        t.Error("Verify .txt should be matched.")
    }
}

func Test_isMatchSuffixes_unmatched(t *testing.T) {
    if isMatchSuffixes([]string{".foo"}, "/tmp/test.txt") {
        t.Error("Verify unmatched suffix")
    }
}

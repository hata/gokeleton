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

func Test_stringsToMap_Nothing(t *testing.T) {
    m := stringsToMap("", ",")
    if len(m) != 0 {
        t.Error("Verify no map key-value found.")
    }
}

func Test_stringsToMap_Simple(t *testing.T) {
    m := stringsToMap("foo=bar", ",")
    if m["foo"] != "bar" {
        t.Error("Verify there is a key value pair.")
    }
}

func Test_stringsToMap_KeyOnly(t *testing.T) {
    m := stringsToMap("foo", ",")
    if m["foo"] != "" {
        t.Error("Verify no value is set")
    }
    m = stringsToMap("foo=", ",")
    if m["foo"] != "" {
        t.Error("Verify no value is set")
    }
}

func Test_stringsToMap_MultiKeywords(t *testing.T) {
    m := stringsToMap("foo=bar,bar=hoge", ",")

    if len(m) != 2 {
        t.Error("Verify there are two entries.")
    }

    if m["foo"] != "bar" {
        t.Error("Verify there is a key value pair.")
    }

    if m["bar"] != "hoge" {
        t.Error("Verify there is a key value pair.")
    }
}

func Test_newSourceAccess_URL(t *testing.T) {
    sa := newSourceAccess("https://github.com/hata/gorep")
    if sa == nil {
        t.Error("Verify https protocol should return SourceAccess for github")
    }
}

func Test_newSourceAccess_File(t *testing.T) {
    sa := newSourceAccess("/tmp")
    if sa == nil {
        t.Error("Verify a local file should return SourceAccess for github")
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

func Test_toSubPath(t *testing.T) {
    path := toSubPath("/tmp/", "/tmp/foo")
    if path != "foo" {
        t.Error("Verify sub path is an expected value")
    }
}

func Test_newReplaceFunc_contents(t *testing.T) {
    rf := newReplaceFunc(map[string]string{"foo":"bar"})
    subPath, contents, err := rf("subPath", "foo,foo,foo,bar,bar,bar")
    if subPath != "subPath" {
        t.Error("Verify subPath is returned.")
    }
    if contents != "bar,bar,bar,bar,bar,bar" {
        t.Error("Verify replacing keywords work well.")
    }
    if err != nil {
        t.Error("Verify no error found")
    }
}


func Test_newReplaceFunc_subPath(t *testing.T) {
    rf := newReplaceFunc(map[string]string{"foo":"bar"})
    subPath, contents, err := rf("foo/foo", "foo,bar,fo")
    if subPath != "bar/bar" {
        t.Error("Verify bar/bar is returned.")
    }
    if contents != "bar,bar,fo" {
        t.Error("Verify replacing keywords work well.")
    }
    if err != nil {
        t.Error("Verify no error found")
    }
}

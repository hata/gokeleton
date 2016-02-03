package main

import (
    "fmt"
    "testing"
)


const sampleURL = "https://github.com/hata/gorep"

func Test_newGithubAccess(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    if ga.client == nil {
        t.Error("Verify github.Client is initialized.")
    }
    if ga.url == "" {
        t.Error("Verify github url is set")
    }
}

func Test_parseURL(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    err := ga.parseURL()
    if err != nil {
        t.Log(err)
        t.Error("error should not be return.")
    }

    if ga.owner != "hata" {
        t.Error("Verify owner should be set")
    }

    if ga.repos != "gorep" {
        t.Error("Verify repos should be set")
    }

    if ga.basePath != "" {
        t.Error("Verify path should be set")
    }
}

func Test_parseURL_directory(t *testing.T) {
    ga := newGithubAccess(sampleURL + "/tree/master/book")
    err := ga.parseURL()
    if err != nil {
        t.Log(err)
        t.Error("error should not be return.")
    }

    if ga.owner != "hata" {
        t.Error("Verify owner should be set")
    }

    if ga.repos != "gorep" {
        t.Error("Verify repos should be set")
    }

    if ga.basePath != "book" {
        t.Error("Verify path should be set")
    }
}

func Test_parseURL_wrong_url(t *testing.T) {
    ga := newGithubAccess("http://www.google.com")
    err := ga.parseURL()
    if err == nil {
        t.Error("Verify error should be return.")
    }
}

func Test_getZipArchive(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    zipReader, err := ga.getZipArchive()
    if err != nil {
        t.Error("There is an error to get archive" + err.Error())
    }
    if zipReader == nil {
        t.Error("No zip body is returned.")
    }
    if len(zipReader.File) == 0 {
        t.Error("No zip len is returned.")
    }
}

func Test_getZipArchive_checkZipArchive(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    zipReader, err := ga.getZipArchive()
    if err != nil {
        t.Error("There is an error to getarchive 2 " + err.Error())
    }

    for _, f := range zipReader.File {
        fmt.Println("File: ", *f)
    }
}

func Test_EachSource(t *testing.T) {
    found := false
    ga := newGithubAccess(sampleURL)
    ga.EachSource(func (fs FileSource) error {
        if fs.SubPath() == ".gitignore" {
            found = true
        }
        fmt.Println("Each File: " + fs.SubPath())
        return nil
    })
    if !found {
        t.Error("Verify there is a file.")
    }
}

func Test_EachSource_for_directory(t *testing.T) {
    gitIgnoreFound := false
    fileFound := false
    ga := newGithubAccess(sampleURL + "/book")
    ga.EachSource(func (fs FileSource) error {
        fmt.Println("Each File: " + fs.SubPath())
        if fs.SubPath() == ".gitignore" {
            gitIgnoreFound = true
        }
        if fs.SubPath() == "chapter.go" {
            fileFound = true
        }
        return nil
    })

    if gitIgnoreFound {
        t.Error("Verify this should not be found")
    }
    if !fileFound {
        t.Error("Verify a file found")
    }
}


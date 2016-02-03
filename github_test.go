package main

import (
    "encoding/base64"
    "fmt"
//    "strings"
    "testing"
)


const sampleURL = "https://github.com/hata/gorep"

func Test_newGithubAccess(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    if ga.client == nil {
        t.Error("Verify github.Client is initialized.")
    }
}

func Test_newGithubAccess_GetReadme(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    repos, _, _ := ga.client.Repositories.GetReadme("hata", "gorep", nil)

    if *repos.Encoding == "base64" {
        _, err := base64.StdEncoding.DecodeString(*repos.Content)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
    }
}

func Test_parseURL(t *testing.T) {
    ga := newGithubAccess(sampleURL)
    err := ga.parseURL(sampleURL)
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
    ga := newGithubAccess(sampleURL)
    err := ga.parseURL(sampleURL + "/tree/master/book")
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
    ga := newGithubAccess(sampleURL)
    err := ga.parseURL("http://www.google.com/")
    if err == nil {
        t.Error("Verify error should be return.")
    }
}
/*
func Test_getReadme(t *testing.T) {
    ga := newGithubAccess()
    readme, err := ga.getReadme("https://github.com/hata/gorep")
    if err != nil {
        t.Error("err should not be return.")
    }
    if readme == "" {
        t.Error("README.md should be return")
    }
    if strings.Index(readme, "gorep") == -1 {
        t.Error("Verify README.md using repository name.")
    }
}

func Test_getContents(t *testing.T) {
    ga := newGithubAccess()
    contents, err := ga.getContents("https://github.com/hata/gorep/README.md")
    if contents == nil && err == nil {
        t.Error("No contents are returned.")
    }
}

*/

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
    ga := newGithubAccess(sampleURL)
    ga.EachSource(func (fs FileSource) error {
        fmt.Println("Each File: " + fs.Name())
        return nil
    })
}


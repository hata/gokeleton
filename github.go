package main

import (
    "archive/zip"
    "bytes"
    "errors"
    "github.com/google/go-github/github"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

type githubAccess struct {
    client *github.Client
    url string
    owner string
    repos string
    basePath string
}

type githubFileSource struct {
    file *zip.File
    path string
}

func newGithubAccess(githubHTMLURL string) (ga *githubAccess) {
    ga = new(githubAccess)
    ga.client = github.NewClient(nil)
    ga.url = githubHTMLURL
    return
}

func newGithubFileSource(zipFile *zip.File, path string) (gf *githubFileSource) {
    gf = new(githubFileSource)
    gf.file = zipFile
    gf.path = path
    return
}

// SourceAccess
func (ga *githubAccess) EachSource(callback FileSourceFunc) (err error) {
    zipReader, err := ga.getZipArchive()
    if err != nil {
        return err
    }

    for _, f := range zipReader.File {
        name := f.Name
        index := strings.Index(name, "/")
        if index != -1 {
            name = name[index + 1:]
        }
        baseLen := len(ga.basePath)
        if baseLen > 0 {
            index = strings.Index(name, ga.basePath)
            if index != 0 {
                continue
            }
            name = name[baseLen + 1:]
        }

        // Check file should be called or not
        err = callback(newGithubFileSource(f, name))
        if err != nil {
            return err
        }
    }

    return nil
}

// FileSource
func (gf *githubFileSource) SubPath() string {
    return gf.path
}

func (gf *githubFileSource) IsDir() bool {
    return gf.file.FileInfo().IsDir()
}

func (gf *githubFileSource) Reader() (io.ReadCloser, error) {
    return gf.file.Open()
}

func (ga *githubAccess) getZipArchive() (zipReader *zip.Reader, err error) {
    var archiveURL *url.URL
    var httpResponse *http.Response
    var zipBytes []byte

    err = ga.parseURL()
    if err != nil {
        return nil, err
    }

    archiveURL, _, err = ga.client.Repositories.GetArchiveLink(ga.owner, ga.repos, github.Zipball, nil)
    if err != nil {
        return nil, err
    }

    httpResponse, err = http.Get(archiveURL.String())
    if err != nil {
        return nil, err
    }

    zipBytes, err = ioutil.ReadAll(httpResponse.Body)
    if err != nil {
        return nil, err
    }

    zipReader, err = zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
    if err != nil {
        return nil, err
    }

    return
}

func (ga *githubAccess) parseURL() error {
    url, err := url.Parse(ga.url)
    if err != nil {
        return err
    }

    if url.Host != "github.com" {
        return errors.New("Host is not matched.")
    }

    pathElements := strings.Split(url.Path, "/")
    if len(pathElements) < 3 || pathElements[1] == "" || pathElements[2] == "" {
        return errors.New("There is no owner and/or repository in url")
    }

    ga.owner = pathElements[1]
    ga.repos = pathElements[2]

    if len(pathElements) < 6 {
        if len(pathElements) == 3 {
            ga.basePath = ""
        } else {
            ga.basePath = strings.Join(pathElements[3:], "/")
        }
        return nil
    }

    if !((pathElements[3] == "tree" || pathElements[4] == "blob") && pathElements[4] == "master") {
        return errors.New("No supported url format. Expected format is (tree|blob)/master.")
    }

    ga.basePath = strings.Join(pathElements[5:], "/")
    return nil
}

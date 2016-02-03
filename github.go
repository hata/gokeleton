package main

import (
    "archive/zip"
    "bytes"
    "encoding/base64"
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
}

func newGithubAccess(githubHTMLURL string) (ga *githubAccess) {
    ga = new(githubAccess)
    ga.client = github.NewClient(nil)
    ga.url = githubHTMLURL
    return
}

func newGithubFileSource(zipFile *zip.File) (gf *githubFileSource) {
    gf = new(githubFileSource)
    gf.file = zipFile
    return
}

// SourceAccess
func (ga *githubAccess) EachSource(callback FileSourceFunc) (err error) {
    zipReader, err := ga.getZipArchive()
    if err != nil {
        return err
    }

    for _, f := range zipReader.File {
        err = callback(newGithubFileSource(f))
        if err != nil {
            return err
        }
    }

    return nil
}

// FileSource
func (gf *githubFileSource) Name() (name string) {
    name = gf.file.Name
    index := strings.Index(name, "/")
    if index != -1 {
        name = name[index + 1:]
    }
    return
}

func (gf *githubFileSource) Reader() (io.ReadCloser, error) {
    return gf.file.Open()
}


/*
func (ga *githubAccess) getReadme(githubURL string) (readme string, err error) {
    err = ga.parseURL(githubURL)
    if err != nil {
        return "", err
    }

    fileContent, _, err := ga.client.Repositories.GetReadme(ga.owner, ga.repos, nil)
    fmt.Println("error: ", err)
    fmt.Println("Encoding: ", *fileContent.Encoding)

    content, err := decodeContentString(fileContent)
    if err != nil {
        return "", err
    }

    return content, nil
}

func (ga *githubAccess) getContents(githubBlobURL string) (contents []byte, err error) {
    err = ga.parseURL(githubBlobURL)
    if err != nil {
        return nil, err
    }

    fileContent, directoryContent, resp, err := ga.client.Repositories.GetContents(ga.owner, ga.repos, ga.basePath, nil)

    fmt.Println("basePath:", ga.basePath)
    fmt.Println("fileContent:", fileContent)
    fmt.Println("directoryContent:", directoryContent)
    fmt.Println("resp:", resp)
    fmt.Println("err:", err)

    if err != nil {
        return
    }

    if fileContent != nil {
        contents, err = decodeContent(fileContent)
    }

    return
}
*/

func (ga *githubAccess) getZipArchive() (zipReader *zip.Reader, err error) {
    var archiveURL *url.URL
    var httpResponse *http.Response
    var zipBytes []byte

    err = ga.parseURL(ga.url)
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

func (ga *githubAccess) parseURL(githubURL string) error {
    url, err := url.Parse(githubURL)
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

/*
func (ga *githubAccess)walkFiles(githubHTMLURL string) (map[string]string, error) {
    err := ga.parseURL(githubHTMLURL)
    if err != nil {
        return nil, err
    }

    contentMap := map[string]string{}

    err = ga.walkContents(contentMap, ga.basePath)
    if err != nil {
        return nil, err
    }

    return contentMap, nil
}

func (ga *githubAccess)walkContents(contentMap map[string]string, path string) error {
    fileContent, directoryContent, _, err := ga.client.Repositories.GetContents(ga.owner, ga.repos, path, nil)

    if err != nil {
        return err
    }

    if directoryContent != nil && fileContent == nil {
        for _, dirContent := range directoryContent {
            ga.walkContents(contentMap, *dirContent.Path)
        }
    } else if fileContent != nil {
        var data string
        data, err = decodeContentString(fileContent)
        if err != nil {
            return err
        }
        contentMap[*fileContent.Path] = data  // *fileContent.DownloadURL
    } else {
        return errors.New("No data was returned.")
    }

    return nil
}
*/

func decodeContentString(repos *github.RepositoryContent) (string, error) {
    data, err := decodeContent(repos)
    if err != nil {
        return "", err
    } else {
        return string(data), nil
    }
}

func decodeContent(repos *github.RepositoryContent) ([]byte, error) {
    if *repos.Encoding == "base64" {
        data, err := base64.StdEncoding.DecodeString(*repos.Content)
        if err != nil {
            return nil, err
        }
        return data, nil
    } else {
        return nil, errors.New("No supported encoding is found")
    }
}


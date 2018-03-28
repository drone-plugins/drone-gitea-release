package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event string
	}

	Commit struct {
		Ref string
	}

	Config struct {
		APIKey     string
		Files      []string
		FileExists string
		Checksum   []string
		Draft      bool
		PreRelease bool
		Insecure   bool
		BaseURL    string
		Title      string
		Note       string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	var (
		files []string
	)

	if p.Build.Event != "tag" {
		return fmt.Errorf("The Gitea Release plugin is only available for tags")
	}

	if p.Config.APIKey == "" {
		return fmt.Errorf("You must provide an API key")
	}

	if !fileExistsValues[p.Config.FileExists] {
		return fmt.Errorf("Invalid value for file_exists")
	}

	if p.Config.BaseURL == "" {
		return fmt.Errorf("You must provide a base url.")
	}

	if !strings.HasSuffix(p.Config.BaseURL, "/") {
		p.Config.BaseURL = p.Config.BaseURL + "/"
	}

	var err error
	if p.Config.Note != "" {
		if p.Config.Note, err = readStringOrFile(p.Config.Note); err != nil {
			return fmt.Errorf("error while reading %s: %v", p.Config.Note, err)
		}
	}

	if p.Config.Title != "" {
		if p.Config.Title, err = readStringOrFile(p.Config.Title); err != nil {
			return fmt.Errorf("error while reading %s: %v", p.Config.Note, err)
		}
	}

	for _, glob := range p.Config.Files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			return fmt.Errorf("Failed to glob %s. %s", glob, err)
		}

		if globed != nil {
			files = append(files, globed...)
		}
	}

	if len(p.Config.Checksum) > 0 {
		var (
			err error
		)

		files, err = writeChecksums(files, p.Config.Checksum)

		if err != nil {
			return fmt.Errorf("Failed to write checksums. %s", err)
		}
	}

	client := gitea.NewClient(p.Config.BaseURL, p.Config.APIKey)

	if p.Config.Insecure {
		cookieJar, _ := cookiejar.New(nil)

		var insecureClient = &http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		client.SetHTTPClient(insecureClient)
	}

	rc := releaseClient{
		Client:     client,
		Owner:      p.Repo.Owner,
		Repo:       p.Repo.Name,
		Tag:        strings.TrimPrefix(p.Commit.Ref, "refs/tags/"),
		Draft:      p.Config.Draft,
		Prerelease: p.Config.PreRelease,
		FileExists: p.Config.FileExists,
		Title:      p.Config.Title,
		Note:       p.Config.Note,
	}

	release, err := rc.buildRelease()

	if err != nil {
		return fmt.Errorf("Failed to create the release. %s", err)
	}

	if err := rc.uploadFiles(release.ID, files); err != nil {
		return fmt.Errorf("Failed to upload the files. %s", err)
	}

	return nil
}

func readStringOrFile(input string) (string, error) {
	// Check if input is a file path
	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		// No file found => use input as result
		return input, nil
	} else if err != nil {
		return "", err
	}
	result, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

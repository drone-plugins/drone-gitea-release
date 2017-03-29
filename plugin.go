package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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
		Prerelease bool
		BaseURL    string
		UploadURL  string
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
		return fmt.Errorf("The GitHub Release plugin is only available for tags")
	}

	if p.Config.APIKey == "" {
		return fmt.Errorf("You must provide an API key")
	}

	if !fileExistsValues[p.Config.FileExists] {
		return fmt.Errorf("Invalid value for file_exists")
	}

	if !strings.HasSuffix(p.Config.BaseURL, "/") {
		p.Config.BaseURL = p.Config.BaseURL + "/"
	}

	if !strings.HasSuffix(p.Config.UploadURL, "/") {
		p.Config.UploadURL = p.Config.UploadURL + "/"
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

	baseURL, err := url.Parse(p.Config.BaseURL)

	if err != nil {
		return fmt.Errorf("Failed to parse base URL. %s", err)
	}

	uploadURL, err := url.Parse(p.Config.UploadURL)

	if err != nil {
		return fmt.Errorf("Failed to parse upload URL. %s", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Config.APIKey})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	client.BaseURL = baseURL
	client.UploadURL = uploadURL

	rc := releaseClient{
		Client:     client,
		Owner:      p.Repo.Owner,
		Repo:       p.Repo.Name,
		Tag:        filepath.Base(p.Commit.Ref),
		Draft:      p.Config.Draft,
		Prerelease: p.Config.Prerelease,
		FileExists: p.Config.FileExists,
	}

	release, err := rc.buildRelease()

	if err != nil {
		return fmt.Errorf("Failed to create the release. %s", err)
	}

	if err := rc.uploadFiles(*release.ID, files); err != nil {
		return fmt.Errorf("Failed to upload the files. %s", err)
	}

	return nil
}

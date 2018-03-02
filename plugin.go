package main

import (
	"fmt"
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
		Prerelease bool
		BaseURL    string
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

	rc := releaseClient{
		Client:     client,
		Owner:      p.Repo.Owner,
		Repo:       p.Repo.Name,
		Tag:        strings.TrimPrefix(p.Commit.Ref, "refs/tags/"),
		Draft:      p.Config.Draft,
		Prerelease: p.Config.Prerelease,
		FileExists: p.Config.FileExists,
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

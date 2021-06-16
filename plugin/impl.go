// Copyright (c) 2021, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// Settings for the plugin.
type Settings struct {
	APIKey     string
	Files      cli.StringSlice
	FileExists string
	Checksum   cli.StringSlice
	Draft      bool
	PreRelease bool
	BaseURL    string
	Title      string
	Note       string

	uploads []string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	if p.pipeline.Build.Event != "tag" {
		return fmt.Errorf("gitea release plugin is only available for tags")
	}

	if p.settings.APIKey == "" {
		return fmt.Errorf("no api key provided")
	}

	if !fileExistsValues[p.settings.FileExists] {
		return fmt.Errorf("invalid value for file_exists")
	}

	if p.settings.BaseURL == "" {
		return fmt.Errorf("no base url provided")
	}

	if !strings.HasSuffix(p.settings.BaseURL, "/") {
		p.settings.BaseURL = p.settings.BaseURL + "/"
	}

	if p.settings.Note != "" {
		if p.settings.Note, err = readStringOrFile(p.settings.Note); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.settings.Note, err)
		}
	}

	if p.settings.Title != "" {
		if p.settings.Title, err = readStringOrFile(p.settings.Title); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.settings.Note, err)
		}
	}

	files := p.settings.Files.Value()
	for _, glob := range files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			return fmt.Errorf("failed to glob %s: %w", glob, err)
		}

		if globed != nil {
			p.settings.uploads = append(p.settings.uploads, globed...)
		}
	}

	if len(files) > 0 && len(p.settings.uploads) < 1 {
		return fmt.Errorf("failed to find any file to release")
	}

	checksum := p.settings.Checksum.Value()
	if len(checksum) > 0 {
		p.settings.uploads, err = writeChecksums(files, checksum)

		if err != nil {
			return fmt.Errorf("failed to write checksums: %w", err)
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	client, err := gitea.NewClient(p.settings.BaseURL, gitea.SetToken(p.settings.APIKey), gitea.SetHTTPClient(p.network.Client))
	if err != nil {
		return err
	}

	rc := releaseClient{
		Client:     client,
		Owner:      p.pipeline.Repo.Owner,
		Repo:       p.pipeline.Repo.Name,
		Tag:        strings.TrimPrefix(p.pipeline.Commit.Ref, "refs/tags/"),
		Draft:      p.settings.Draft,
		Prerelease: p.settings.PreRelease,
		FileExists: p.settings.FileExists,
		Title:      p.settings.Title,
		Note:       p.settings.Note,
	}

	// When the title was not provided in the config use the tag instead
	if rc.Title == "" {
		rc.Title = rc.Tag
	}

	release, err := rc.buildRelease()

	if err != nil {
		return fmt.Errorf("failed to create the release. %s", err)
	}

	if err := rc.uploadFiles(release.ID, p.settings.uploads); err != nil {
		return fmt.Errorf("failed to upload the files. %s", err)
	}

	return nil
}

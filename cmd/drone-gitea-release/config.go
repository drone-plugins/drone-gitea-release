// Copyright (c) 2021, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-gitea-release/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{

		&cli.StringFlag{
			Name:        "api-key",
			Usage:       "api key to access gitea api",
			EnvVars:     []string{"PLUGIN_API_KEY", "GITEA_RELEASE_API_KEY", "GITEA_TOKEN"},
			Destination: &settings.APIKey,
		},
		&cli.StringSliceFlag{
			Name:        "files",
			Usage:       "list of files to upload",
			EnvVars:     []string{"PLUGIN_FILES", "GITEA_RELEASE_FILES"},
			Destination: &settings.Files,
		},
		&cli.StringFlag{
			Name:        "file-exists",
			Value:       "overwrite",
			Usage:       "what to do if file already exist",
			EnvVars:     []string{"PLUGIN_FILE_EXISTS", "GITEA_RELEASE_FILE_EXISTS"},
			Destination: &settings.FileExists,
		},
		&cli.StringSliceFlag{
			Name:        "checksum",
			Usage:       "generate specific checksums",
			EnvVars:     []string{"PLUGIN_CHECKSUM", "GITEA_RELEASE_CHECKSUM"},
			Destination: &settings.Checksum},
		&cli.BoolFlag{
			Name:        "draft",
			Usage:       "create a draft release",
			EnvVars:     []string{"PLUGIN_DRAFT", "GITEA_RELEASE_DRAFT"},
			Destination: &settings.Draft},
		&cli.BoolFlag{
			Name:        "prerelease",
			Usage:       "set the release as prerelease",
			EnvVars:     []string{"PLUGIN_PRERELEASE", "GITEA_RELEASE_PRERELEASE"},
			Destination: &settings.PreRelease},
		&cli.StringFlag{
			Name:        "base-url",
			Usage:       "url of the gitea instance",
			EnvVars:     []string{"PLUGIN_BASE_URL", "GITEA_RELEASE_BASE_URL"},
			Destination: &settings.BaseURL},
		&cli.StringFlag{
			Name:        "title",
			Value:       "",
			Usage:       "file or string for the title shown in the gitea release",
			EnvVars:     []string{"PLUGIN_TITLE", "GITEA_RELEASE_TITLE"},
			Destination: &settings.Title},
		&cli.StringFlag{
			Name:        "note",
			Value:       "",
			Usage:       "file or string with notes for the release (example: changelog)",
			EnvVars:     []string{"PLUGIN_NOTE", "GITEA_RELEASE_NOTE"},
			Destination: &settings.Note},
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "github-release plugin"
	app.Usage = "github-release plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "api key to access github api",
			EnvVar: "PLUGIN_API_KEY,GITHUB_RELEASE_API_KEY",
		},
		cli.StringSliceFlag{
			Name:   "files",
			Usage:  "list of files to upload",
			EnvVar: "PLUGIN_FILES,GITHUB_RELEASE_FILES",
		},
		cli.StringFlag{
			Name:   "file-exists",
			Value:  "overwrite",
			Usage:  "what to do if file already exist",
			EnvVar: "PLUGIN_FILE_EXISTS,GITHUB_RELEASE_FILE_EXISTS",
		},
		cli.StringSliceFlag{
			Name:   "checksum",
			Usage:  "generate specific checksums",
			EnvVar: "PLUGIN_CHECKSUM,GITHUB_RELEASE_CHECKSUM",
		},
		cli.BoolFlag{
			Name:   "draft",
			Usage:  "create a draft release",
			EnvVar: "PLUGIN_DRAFT,GITHUB_RELEASE_DRAFT",
		},
		cli.StringFlag{
			Name:   "base-url",
			Value:  "https://api.github.com/",
			Usage:  "api url, needs to be changed for ghe",
			EnvVar: "PLUGIN_BASE_URL,GITHUB_RELEASE_BASE_URL",
		},
		cli.StringFlag{
			Name:   "upload-url",
			Value:  "https://uploads.github.com/",
			Usage:  "upload url, needs to be changed for ghe",
			EnvVar: "PLUGIN_UPLOAD_URL,GITHUB_RELEASE_UPLOAD_URL",
		},

		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Event: c.String("build.event"),
		},
		Commit: Commit{
			Ref: c.String("commit.sha"),
		},
		Config: Config{
			APIKey:     c.String("api-key"),
			Files:      c.StringSlice("files"),
			FileExists: c.String("file-exists"),
			Checksum:   c.StringSlice("checksum"),
			Draft:      c.Bool("draft"),
			BaseURL:    c.String("base-url"),
			UploadURL:  c.String("upload-url"),
		},
	}

	return plugin.Exec()
}

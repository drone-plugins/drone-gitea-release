package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitea-release plugin"
	app.Usage = "gitea-release plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "api key to access gitea api",
			EnvVar: "PLUGIN_API_KEY,GITEA_RELEASE_API_KEY,GITEA_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "files",
			Usage:  "list of files to upload",
			EnvVar: "PLUGIN_FILES,GITEA_RELEASE_FILES",
		},
		cli.StringFlag{
			Name:   "file-exists",
			Value:  "overwrite",
			Usage:  "what to do if file already exist",
			EnvVar: "PLUGIN_FILE_EXISTS,GITEA_RELEASE_FILE_EXISTS",
		},
		cli.StringSliceFlag{
			Name:   "checksum",
			Usage:  "generate specific checksums",
			EnvVar: "PLUGIN_CHECKSUM,GITEA_RELEASE_CHECKSUM",
		},
		cli.BoolFlag{
			Name:   "draft",
			Usage:  "create a draft release",
			EnvVar: "PLUGIN_DRAFT,GITEA_RELEASE_DRAFT",
		},
		cli.BoolFlag{
			Name:   "insecure",
			Usage:  "visit base-url via insecure https protocol",
			EnvVar: "PLUGIN_INSECURE,GITEA_RELEASE_INSECURE",
		},
		cli.BoolFlag{
			Name:   "prerelease",
			Usage:  "set the release as prerelease",
			EnvVar: "PLUGIN_PRERELEASE,GITEA_RELEASE_PRERELEASE",
		},
		cli.StringFlag{
			Name:   "base-url",
			Usage:  "url of the gitea instance",
			EnvVar: "PLUGIN_BASE_URL,GITEA_RELEASE_BASE_URL",
		},
		cli.StringFlag{
			Name:   "note",
			Value:  "",
			Usage:  "file or string with notes for the release (example: changelog)",
			EnvVar: "PLUGIN_NOTE,GITEA_RELEASE_NOTE",
		},
		cli.StringFlag{
			Name:   "title",
			Value:  "",
			Usage:  "file or string for the title shown in the gitea release",
			EnvVar: "PLUGIN_TITLE,GITEA_RELEASE_TITLE",
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
			Ref: c.String("commit.ref"),
		},
		Config: Config{
			APIKey:     c.String("api-key"),
			Files:      c.StringSlice("files"),
			FileExists: c.String("file-exists"),
			Checksum:   c.StringSlice("checksum"),
			Draft:      c.Bool("draft"),
			PreRelease: c.Bool("prerelease"),
			BaseURL:    c.String("base-url"),
			Insecure:   c.Bool("insecure"),
			Title:      c.String("title"),
			Note:       c.String("note"),
		},
	}

	return plugin.Exec()
}

package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	build     string
	buildDate string
)

func main() {
	fmt.Printf("Drone GitHub Release Plugin built at %s\n", buildDate)

	workspace := drone.Workspace{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if build.Event != "tag" {
		fmt.Printf("The GitHub Release plugin is only available for tags\n")
		os.Exit(0)

		return
	}

	if len(vargs.BaseUrl) == 0 {
		vargs.BaseUrl = "https://api.github.com/"
	} else {
		if !strings.HasSuffix(vargs.BaseUrl, "/") {
			vargs.BaseUrl = vargs.BaseUrl + "/"
		}
	}

	if len(vargs.UploadUrl) == 0 {
		vargs.UploadUrl = "https://uploads.github.com/"
	} else {
		if !strings.HasSuffix(vargs.UploadUrl, "/") {
			vargs.UploadUrl = vargs.UploadUrl + "/"
		}
	}

	if len(vargs.APIKey) == 0 {
		fmt.Printf("You must provide an API key\n")
		os.Exit(1)

		return
	}

	if len(workspace.Path) != 0 {
		os.Chdir(workspace.Path)
	}

	files := make([]string, 0)

	for _, glob := range vargs.Files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			fmt.Printf("Failed to glob %s\n", glob)
			os.Exit(1)

			return
		}

		if globed != nil {
			files = append(files, globed...)
		}
	}

	baseUrl, err := url.Parse(vargs.BaseUrl)

	if err != nil {
		fmt.Printf("Failed to parse base URL\n")
		os.Exit(1)

		return
	}

	uploadUrl, err := url.Parse(vargs.UploadUrl)

	if err != nil {
		fmt.Printf("Failed to parse upload URL\n")
		os.Exit(1)

		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: vargs.APIKey,
		})

	tc := oauth2.NewClient(
		oauth2.NoContext,
		ts)

	client := github.NewClient(tc)
	client.BaseURL = baseUrl
	client.UploadURL = uploadUrl

	release, releaseErr := prepareRelease(
		client,
		repo.Owner,
		repo.Name,
		filepath.Base(build.Ref))

	if releaseErr != nil {
		fmt.Println(releaseErr)
		os.Exit(1)

		return
	}

	uploadErr := appendFiles(
		client,
		repo.Owner,
		repo.Name,
		*release.ID,
		files)

	if uploadErr != nil {
		fmt.Println(uploadErr)
		os.Exit(1)

		return
	}
}

func prepareRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	release, _ := retrieveRelease(
		client,
		owner,
		repo,
		tag)

	if release != nil {
		return release, nil
	}

	release, err := createRelease(
		client,
		owner,
		repo,
		tag)

	if err != nil {
		return nil, err
	}

	if release != nil {
		return release, nil
	}

	return nil, errors.New(
		"Failed to retrieve or create a release")
}

func retrieveRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	release, _, err := client.Repositories.GetReleaseByTag(
		owner,
		repo,
		tag)

	if err != nil {
		return nil, errors.New(
			"Failed to retrieve release")
	}

	fmt.Printf("Successfully retrieved %s release\n", tag)
	return release, nil
}

func createRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	release, _, err := client.Repositories.CreateRelease(
		owner,
		repo,
		&github.RepositoryRelease{TagName: github.String(tag)})

	if err != nil {
		return nil, errors.New(
			"Failed to create release")
	}

	fmt.Printf("Successfully created %s release\n", tag)
	return release, nil
}

func appendFiles(client *github.Client, owner string, repo string, id int, files []string) error {
	assets, _, err := client.Repositories.ListReleaseAssets(
		owner,
		repo,
		id,
		&github.ListOptions{})

	if err != nil {
		return errors.New(
			"Failed to fetch existing assets")
	}

	for _, file := range files {
		handle, err := os.Open(file)

		if err != nil {
			return errors.New(
				fmt.Sprintf("Failed to read %s artifact", file))
		}

		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				_, deleteErr := client.Repositories.DeleteReleaseAsset(
					owner,
					repo,
					*asset.ID)

				if deleteErr != nil {
					return errors.New(
						fmt.Sprintf("Failed to delete %s artifact", file))
				} else {
					fmt.Printf("Successfully deleted old %s artifact\n", *asset.Name)
				}
			}
		}

		_, _, uploadErr := client.Repositories.UploadReleaseAsset(
			owner,
			repo,
			id,
			&github.UploadOptions{Name: path.Base(file)},
			handle)

		if uploadErr != nil {
			return errors.New(
				fmt.Sprintf("Failed to upload %s artifact", file))
		} else {
			fmt.Printf("Successfully uploaded %s artifact\n", file)
		}
	}

	return nil
}

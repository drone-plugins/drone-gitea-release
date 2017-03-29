package main

import (
	"fmt"
	"os"
	"path"

	"github.com/google/go-github/github"
)

// Release holds ties the drone env data and github client together.
type releaseClient struct {
	*github.Client
	Owner      string
	Repo       string
	Tag        string
	Draft      bool
	Prerelease bool
	FileExists string
}

func (rc *releaseClient) buildRelease() (*github.RepositoryRelease, error) {
	// first attempt to get a release by that tag
	release, err := rc.getRelease()

	if err != nil && release == nil {
		fmt.Println(err)
	} else if release != nil {
		return release, nil
	}

	// if no release was found by that tag, create a new one
	release, err = rc.newRelease()

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve or create a release: %s", err)
	}

	return release, nil
}

func (rc *releaseClient) getRelease() (*github.RepositoryRelease, error) {
	release, _, err := rc.Client.Repositories.GetReleaseByTag(rc.Owner, rc.Repo, rc.Tag)

	if err != nil {
		return nil, fmt.Errorf("Release %s not found", rc.Tag)
	}

	fmt.Printf("Successfully retrieved %s release\n", rc.Tag)
	return release, nil
}

func (rc *releaseClient) newRelease() (*github.RepositoryRelease, error) {
	rr := &github.RepositoryRelease{
		TagName:    github.String(rc.Tag),
		Draft:      &rc.Draft,
		Prerelease: &rc.Prerelease,
	}

	release, _, err := rc.Client.Repositories.CreateRelease(rc.Owner, rc.Repo, rr)

	if err != nil {
		return nil, fmt.Errorf("Failed to create release: %s", err)
	}

	fmt.Printf("Successfully created %s release\n", rc.Tag)
	return release, nil
}

func (rc *releaseClient) uploadFiles(id int, files []string) error {
	assets, _, err := rc.Client.Repositories.ListReleaseAssets(rc.Owner, rc.Repo, id, &github.ListOptions{})

	if err != nil {
		return fmt.Errorf("Failed to fetch existing assets: %s", err)
	}

	var uploadFiles []string

files:
	for _, file := range files {
		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				switch rc.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("Asset file %s already exists", path.Base(file))
				case "skip":
					fmt.Printf("Skipping pre-existing %s artifact\n", *asset.Name)
					continue files
				default:
					return fmt.Errorf("Internal error, unkown file_exist value %s", rc.FileExists)
				}
			}
		}

		uploadFiles = append(uploadFiles, file)
	}

	for _, file := range uploadFiles {
		handle, err := os.Open(file)

		if err != nil {
			return fmt.Errorf("Failed to read %s artifact: %s", file, err)
		}

		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				if _, err := rc.Client.Repositories.DeleteReleaseAsset(rc.Owner, rc.Repo, *asset.ID); err != nil {
					return fmt.Errorf("Failed to delete %s artifact: %s", file, err)
				}

				fmt.Printf("Successfully deleted old %s artifact\n", *asset.Name)
			}
		}

		uo := &github.UploadOptions{Name: path.Base(file)}

		if _, _, err = rc.Client.Repositories.UploadReleaseAsset(rc.Owner, rc.Repo, id, uo, handle); err != nil {
			return fmt.Errorf("Failed to upload %s artifact: %s", file, err)
		}

		fmt.Printf("Successfully uploaded %s artifact\n", file)
	}

	return nil
}

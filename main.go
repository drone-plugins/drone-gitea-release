package main

import (
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
	}

	if vargs.BaseURL == "" {
		vargs.BaseURL = "https://api.github.com/"
	} else if !strings.HasSuffix(vargs.BaseURL, "/") {
		vargs.BaseURL = vargs.BaseURL + "/"
	}

	if vargs.UploadURL == "" {
		vargs.UploadURL = "https://uploads.github.com/"
	} else if !strings.HasSuffix(vargs.UploadURL, "/") {
		vargs.UploadURL = vargs.UploadURL + "/"
	}

	if vargs.APIKey == "" {
		fmt.Printf("You must provide an API key\n")
		os.Exit(1)
	}

	if workspace.Path != "" {
		os.Chdir(workspace.Path)
	}

	var files []string
	for _, glob := range vargs.Files.Slice() {
		globed, err := filepath.Glob(glob)
		if err != nil {
			fmt.Printf("Failed to glob %s\n", glob)
			os.Exit(1)
		}
		if globed != nil {
			files = append(files, globed...)
		}
	}

	if vargs.Checksum.Len() > 0 {
		var err error
		files, err = writeChecksums(files, vargs.Checksum.Slice())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	baseURL, err := url.Parse(vargs.BaseURL)
	if err != nil {
		fmt.Printf("Failed to parse base URL\n")
		os.Exit(1)
	}

	uploadURL, err := url.Parse(vargs.UploadURL)
	if err != nil {
		fmt.Printf("Failed to parse upload URL\n")
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: vargs.APIKey})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	client.BaseURL = baseURL
	client.UploadURL = uploadURL

	rc := releaseClient{
		Client: client,
		Owner:  repo.Owner,
		Repo:   repo.Name,
		Tag:    filepath.Base(build.Ref),
		Draft:  vargs.Draft,
	}

	release, err := rc.buildRelease()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rc.uploadFiles(*release.ID, files); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Release holds ties the drone env data and github client together.
type releaseClient struct {
	*github.Client
	Owner string
	Repo  string
	Tag   string
	Draft bool
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
		TagName: github.String(rc.Tag),
		Draft: &rc.Draft,
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

	for _, file := range files {
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

func writeChecksums(files, methods []string) ([]string, error) {

	checksums := make(map[string][]string)
	for _, method := range methods {
		for _, file := range files {
			handle, err := os.Open(file)
			if err != nil {
				return nil, fmt.Errorf("Failed to read %s artifact: %s", file, err)
			}

			hash, err := checksum(handle, method)
			if err != nil {
				return nil, err
			}

			checksums[method] = append(checksums[method], hash, file)
		}
	}

	for method, results := range checksums {
		filename := method + "sum.txt"
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(results); i += 2 {
			hash := results[i]
			file := results[i+1]
			if _, err := f.WriteString(fmt.Sprintf("%s  %s\n", hash, file)); err != nil {
				return nil, err
			}
		}
		files = append(files, filename)
	}
	return files, nil
}

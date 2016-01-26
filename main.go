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

	release, err := buildRelease(client, repo.Owner, repo.Name, filepath.Base(build.Ref))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := uploadFiles(client, repo.Owner, repo.Name, *release.ID, files); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {

	// first attempt to get a release by that tag
	release, err := getRelease(client, owner, repo, tag)
	if err != nil && release == nil {
		fmt.Println(err)
	} else if release != nil {
		return release, nil
	}

	// if no release was found by that tag, create a new one
	release, err = newRelease(client, owner, repo, tag)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve or create a release: %s", err)
	}

	return release, nil
}

func getRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	release, _, err := client.Repositories.GetReleaseByTag(owner, repo, tag)
	if err != nil {
		return nil, fmt.Errorf("Release %s not found", tag)
	}

	fmt.Printf("Successfully retrieved %s release\n", tag)
	return release, nil
}

func newRelease(client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	rr := &github.RepositoryRelease{TagName: github.String(tag)}
	release, _, err := client.Repositories.CreateRelease(owner, repo, rr)
	if err != nil {
		return nil, fmt.Errorf("Failed to create release: %s", err)
	}

	fmt.Printf("Successfully created %s release\n", tag)
	return release, nil
}

func uploadFiles(client *github.Client, owner string, repo string, id int, files []string) error {
	assets, _, err := client.Repositories.ListReleaseAssets(owner, repo, id, &github.ListOptions{})
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
				if _, err := client.Repositories.DeleteReleaseAsset(owner, repo, *asset.ID); err != nil {
					return fmt.Errorf("Failed to delete %s artifact: %s", file, err)
				}
				fmt.Printf("Successfully deleted old %s artifact\n", *asset.Name)
			}
		}

		uo := &github.UploadOptions{Name: path.Base(file)}
		if _, _, err = client.Repositories.UploadReleaseAsset(owner, repo, id, uo, handle); err != nil {
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

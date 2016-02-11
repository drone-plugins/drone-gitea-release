package main

import "github.com/drone/drone-go/drone"

// Params are the parameters that the GitHub Release plugin can parse.
type Params struct {
	BaseURL   string            `json:"base_url"`
	UploadURL string            `json:"upload_url"`
	APIKey    string            `json:"api_key"`
	Files     drone.StringSlice `json:"files"`
	Checksum  drone.StringSlice `json:"checksum"`
	Draft     bool              `json:"draft"`
}

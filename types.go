package main

type Params struct {
	BaseUrl   string   `json:"base_url"`
	UploadUrl string   `json:"upload_url"`
	APIKey    string   `json:"api_key"`
	Files     []string `json:"files"`
}

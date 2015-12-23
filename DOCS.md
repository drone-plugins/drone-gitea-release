Use this  plugin for publishing files and artifacts to GitHub releases. You
can override the default configuration with the following parameters:

* `api_key` - GitHub oauth token with public_repo or repo permission
* `files` - Files to upload to GitHub Release, globs are allowed
* `base_url` - GitHub base URL, only required for GHE
* `upload_url` - GitHub upload URL, only required for GHE

Sample configuration:

```yaml
publish:
  github_release:
    api_key: my_github_api_key
    files:
      - dist/*
```

Use this  plugin for publishing files and artifacts to GitHub releases. You
can override the default configuration with the following parameters:

* `api_key` - GitHub oauth token with public_repo or repo permission
* `files` - Files to upload to GitHub Release, globs are allowed
* `checksum` - Checksum takes hash methods to include in your GitHub release for the files specified. Supported hash methods include md5, sha1, sha256, sha512, adler32, and crc32.
* `base_url` - GitHub base URL, only required for GHE
* `upload_url` - GitHub upload URL, only required for GHE

Sample configuration:

```yaml
publish:
  github_release:
    api_key: my_github_api_key
    files:
      - dist/*
    checksum: sha1
```

or

```yaml
publish:
  github_release:
    api_key: my_github_api_key
    files:
      - dist/*
    checksum:
      - md5
      - sha1
      - sha256
      - sha512
      - adler32
      - crc32
```

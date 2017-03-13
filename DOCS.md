Use this plugin for publishing files and artifacts to GitHub releases. Be aware
that you can use this plugin only for tags, GitHub doesn't support the release
of branches.

## Config

The following parameters are used to configure the plugin:

* **api_key** - GitHub oauth token with public_repo or repo permission
* **files** - files to upload to GitHub Release, globs are allowed
* **file_exists** - what to do if an file asset already exists, supported values: **overwrite** (default), **skip** and **fail**
* **checksum** - checksum takes hash methods to include in your GitHub release for the files specified. Supported hash methods include md5, sha1, sha256, sha512, adler32, and crc32.
* **draft** - create a draft release if set to true
* **prerelease** - set the release as prerelease if set to true
* **base_url** - GitHub base URL, only required for GHE
* **upload_url** - GitHub upload URL, only required for GHE

The following secret values can be set to configure the plugin.

* **GITHUB_RELEASE_API_KEY** - corresponds to **api_key**
* **GITHUB_RELEASE_BASE_URL** - corresponds to **base_url**
* **GITHUB_RELEASE_UPLOAD_URL** - corresponds to **upload_url**

It is highly recommended to put the **GITHUB_RELEASE_API_KEY** into a secret so
it is not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=plugins/github-release \
    octocat/hello-world GITHUB_RELEASE_API_KEY my_github_api_key
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Examples

The following is a sample configuration in your .drone.yml file:

```yaml
pipeline:
  github_release:
    image: plugins/github-release
    files: dist/*
    when:
      event: tag
```

An example for generating checksums and upload additional files:

```yaml
pipeline:
  github_release:
    image: plugins/github-release
    files:
      - dist/*
      - bin/binary.exe
    checksum:
      - md5
      - sha1
      - sha256
      - sha512
      - adler32
      - crc32
    when:
      event: tag
```

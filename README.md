# drone-github-release

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-github-release/status.svg)](http://beta.drone.io/drone-plugins/drone-github-release)
[![](https://badge.imagelayers.io/plugins/drone-github-release:latest.svg)](https://imagelayers.io/?images=plugins/drone-github-release:latest 'Get your own badge on imagelayers.io')

Drone plugin for publishing GitHub releases

## Usage

```
./drone-github-release <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "full_name": "drone/drone",
        "owner": "drone",
        "name": "drone"
    },
    "build": {
        "event": "tag",
        "branch": "refs/heads/v0.0.1",
        "commit": "8f5d3b2ce38562bedb48b798328f5bb2e4077a2f",
        "ref": "refs/heads/v0.0.1"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "api_key": "your_api_key",
        "files": [
            "dist/*.txt",
            "dist/other-file"
        ],
        "checksum": [
            "md5",
            "sha1",
            "sha256",
            "sha512",
            "adler32",
            "crc32"
        ]
    }
}
EOF
```

## Docker

Build the Docker container using `make`:

```
make deps build docker
```

### Example

```sh
docker run -i plugins/drone-github-release <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "full_name": "drone/drone",
        "owner": "drone",
        "name": "drone"
    },
    "build": {
        "event": "tag",
        "branch": "refs/heads/v0.0.1",
        "commit": "8f5d3b2ce38562bedb48b798328f5bb2e4077a2f",
        "ref": "refs/heads/v0.0.1"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "api_key": "your_api_key",
        "files": [
            "dist/*.txt",
            "dist/other-file"
        ],
        "checksum": [
            "md5",
            "sha1",
            "sha256",
            "sha512",
            "adler32",
            "crc32"
        ]
    }
}
EOF
```

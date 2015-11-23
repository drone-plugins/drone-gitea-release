# drone-github-release

Drone plugin for publishing GitHub releases

## Usage

Publish a release:

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
        "path": "drone/src/github.com/drone/drone"
    },
    "vargs": {
        "api_key": "your_api_key",
        "files": [
            "dist/*.txt",
            "dist/other-file"
        ]
    }
}
EOF
```

## Docker

Build the Docker container using the `netgo` build tag to eliminate the CGO
dependency:

```
GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 go build -a -tags netgo
docker build --rm=true -t plugins/drone-github-release .
```

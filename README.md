# drone-github-release

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-github-release/status.svg)](http://beta.drone.io/drone-plugins/drone-github-release)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-github-release?status.svg)](http://godoc.org/github.com/drone-plugins/drone-github-release)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-github-release)](https://goreportcard.com/report/github.com/drone-plugins/drone-github-release)
[![](https://images.microbadger.com/badges/image/plugins/github-release.svg)](https://microbadger.com/images/plugins/github-release "Get your own image badge on microbadger.com")

Drone plugin to publish files and artifacts to GitHub Release. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-github-release/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-github-release
docker build --rm -t plugins/github-release .
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e DRONE_BUILD_EVENT=tag \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=foo \
  -e DRONE_COMMIT_REF=refs/heads/master \
  -e PLUGIN_API_KEY=${HOME}/.ssh/id_rsa \
  -e PLUGIN_FILES=master \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/github-release
```

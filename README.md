# drone-github-release

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-github-release/status.svg)](http://beta.drone.io/drone-plugins/drone-github-release)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-github-release?status.svg)](http://godoc.org/github.com/drone-plugins/drone-github-release)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-github-release)](https://goreportcard.com/report/github.com/drone-plugins/drone-github-release)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to publish files and artifacts to GitHub Release. For the usage
information and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/github-release .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-github-release' not found or does not exist..
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

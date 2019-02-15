# drone-gitea-release

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-gitea-release/status.svg)](http://cloud.drone.io/drone-plugins/drone-gitea-release)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/gitea-release.svg)](https://microbadger.com/images/plugins/gitea-release "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-gitea-release?status.svg)](http://godoc.org/github.com/drone-plugins/drone-gitea-release)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-gitea-release)](https://goreportcard.com/report/github.com/drone-plugins/drone-gitea-release)

Drone plugin to publish files and artifacts to Gitea release. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-gitea-release/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-gitea-release
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/gitea-release .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_BASE_URL=https://try.gitea.io \
  -e PLUGIN_API_KEY=your-api-key \
  -e PLUGIN_FILES=build/* \
  -e DRONE_REPO_OWNER=gitea \
  -e DRONE_REPO_NAME=test \
  -e DRONE_BUILD_EVENT=tag \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/gitea-release
```

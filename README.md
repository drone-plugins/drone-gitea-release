# drone-gitea-release

[![Build Status](https://beta.drone.io/api/badges/drone-plugins/drone-gitea-release/status.svg)](https://beta.drone.io/drone-plugins/drone-gitea-release)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-gitea-release?status.svg)](http://godoc.org/github.com/drone-plugins/drone-gitea-release)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-gitea-release)](https://goreportcard.com/report/github.com/drone-plugins/drone-gitea-release)

Drone plugin to publish files and artifacts to Gitea Release.

## Build

Build the binary with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-gitea-release
docker build --rm -t plugins/gitea-release .
```
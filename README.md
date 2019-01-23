# drone-gitea-release

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-gitea-release/status.svg)](http://cloud.drone.io/drone-plugins/drone-gitea-release)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/gitea-release.svg)](https://microbadger.com/images/plugins/gitea-release "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-gitea-release?status.svg)](http://godoc.org/github.com/drone-plugins/drone-gitea-release)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-gitea-release)](https://goreportcard.com/report/github.com/drone-plugins/drone-gitea-release)

Drone plugin to publish files and artifacts to Gitea Release.

**Note: This plugin requires Gitea 1.5 or newer.**

## Build

Build the binary with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-gitea-release
docker build --rm -t plugins/gitea-release .
```

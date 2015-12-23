# Docker image for the Drone GitHub Release plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-github-release
#     make deps build docker

FROM alpine:3.2

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-github-release /bin/
ENTRYPOINT ["/bin/drone-github-release"]

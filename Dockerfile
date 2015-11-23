# Docker image for the Drone GitHub Release plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-github-release
#     GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-github-release .

FROM alpine:3.2
RUN apk add -U ca-certificates && rm -rf /var/cache/apk/*
ADD drone-github-release /bin/
ENTRYPOINT ["/bin/drone-github-release"]

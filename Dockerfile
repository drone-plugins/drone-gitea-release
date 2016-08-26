FROM alpine:3.4

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-github-release /bin/
ENTRYPOINT ["/bin/drone-github-release"]

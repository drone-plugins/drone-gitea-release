# escape=`
FROM plugins/base:windows-1809

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" `
  org.label-schema.name="Drone Gitea Release" `
  org.label-schema.vendor="Drone.IO Community" `
  org.label-schema.schema-version="1.0"

ADD release/windows/amd64/drone-gitea-release.exe C:/bin/drone-gitea-release.exe
ENTRYPOINT [ "C:\\bin\\drone-gitea-release.exe" ]

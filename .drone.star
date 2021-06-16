def main(ctx):
    before = testing(ctx)

    stages = [
        linux(ctx, "amd64"),
        linux(ctx, "arm64"),
        linux(ctx, "arm"),
        windows(ctx, "1909"),
        windows(ctx, "1903"),
        windows(ctx, "1809"),
    ]

    after = manifest(ctx) + gitter(ctx)

    for b in before:
        for s in stages:
            s["depends_on"].append(b["name"])

    for s in stages:
        for a in after:
            a["depends_on"].append(s["name"])

    return before + stages + after

def testing(ctx):
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "testing",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "staticcheck",
                "image": "golang:1.15",
                "pull": "always",
                "commands": [
                    "go run honnef.co/go/tools/cmd/staticcheck ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "lint",
                "image": "golang:1.15",
                "commands": [
                    "go run golang.org/x/lint/golint -set_exit_status ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "vet",
                "image": "golang:1.15",
                "commands": [
                    "go vet ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "test",
                "image": "golang:1.15",
                "commands": [
                    "go test -cover ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
        ],
        "volumes": [
            {
                "name": "gopath",
                "temp": {},
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }]

def linux(ctx, arch):
    if ctx.build.event == "tag":
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-gitea-release ./cmd/drone-gitea-release' % (ctx.build.ref.replace("refs/tags/v", ""), arch),
        ]
    else:
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-gitea-release ./cmd/drone-gitea-release' % (ctx.build.commit[0:8], arch),
        ]

    steps = [
        {
            "name": "environment",
            "image": "golang:1.15",
            "pull": "always",
            "environment": {
                "CGO_ENABLED": "0",
            },
            "commands": [
                "go version",
                "go env",
            ],
        },
        {
            "name": "build",
            "image": "golang:1.15",
            "environment": {
                "CGO_ENABLED": "0",
            },
            "commands": build,
        },
        {
            "name": "executable",
            "image": "golang:1.15",
            "commands": [
                "./release/linux/%s/drone-gitea-release --help" % (arch),
            ],
        },
    ]

    if ctx.build.event != "pull_request":
        steps.append({
            "name": "docker",
            "image": "plugins/docker",
            "settings": {
                "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                "repo": "plugins/gitea-release",
                "username": {
                    "from_secret": "docker_username",
                },
                "password": {
                    "from_secret": "docker_password",
                },
                "auto_tag": True,
                "auto_tag_suffix": "linux-%s" % (arch),
            },
        })

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "linux-%s" % (arch),
        "platform": {
            "os": "linux",
            "arch": arch,
        },
        "steps": steps,
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }

def windows(ctx, version):
    docker = [
        "echo $env:PASSWORD | docker login --username $env:USERNAME --password-stdin",
    ]

    if ctx.build.event == "tag":
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/windows/amd64/drone-gitea-release.exe ./cmd/drone-gitea-release' % (ctx.build.ref.replace("refs/tags/v", "")),
        ]

        docker = docker + [
            "docker build --pull -f docker/Dockerfile.windows.%s -t plugins/gitea-release:%s-windows-%s-amd64 ." % (version, ctx.build.ref.replace("refs/tags/v", ""), version),
            "docker run --rm plugins/gitea-release:%s-windows-%s-amd64 --help" % (ctx.build.ref.replace("refs/tags/v", ""), version),
            "docker push plugins/gitea-release:%s-windows-%s-amd64" % (ctx.build.ref.replace("refs/tags/v", ""), version),
        ]
    else:
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/windows/amd64/drone-gitea-release.exe ./cmd/drone-gitea-release' % (ctx.build.commit[0:8]),
        ]

        docker = docker + [
            "docker build --pull -f docker/Dockerfile.windows.%s -t plugins/gitea-release:windows-%s-amd64 ." % (version, version),
            "docker run --rm plugins/gitea-release:windows-%s-amd64 --help" % (version),
            "docker push plugins/gitea-release:windows-%s-amd64" % (version),
        ]

    return {
        "kind": "pipeline",
        "type": "ssh",
        "name": "windows-%s" % (version),
        "platform": {
            "os": "windows",
        },
        "server": {
            "host": {
                "from_secret": "windows_server_%s" % (version),
            },
            "user": {
                "from_secret": "windows_username",
            },
            "password": {
                "from_secret": "windows_password",
            },
        },
        "steps": [
            {
                "name": "environment",
                "environment": {
                    "CGO_ENABLED": "0",
                },
                "commands": [
                    "go version",
                    "go env",
                ],
            },
            {
                "name": "build",
                "environment": {
                    "CGO_ENABLED": "0",
                },
                "commands": build,
            },
            {
                "name": "executable",
                "commands": [
                    "./release/windows/amd64/drone-gitea-release.exe --help",
                ],
            },
            {
                "name": "docker",
                "environment": {
                    "USERNAME": {
                        "from_secret": "docker_username",
                    },
                    "PASSWORD": {
                        "from_secret": "docker_password",
                    },
                },
                "commands": docker,
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
            ],
        },
    }

def manifest(ctx):
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "manifest",
        "steps": [
            {
                "name": "manifest",
                "image": "plugins/manifest",
                "settings": {
                    "auto_tag": "true",
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "spec": "docker/manifest.tmpl",
                    "ignore_missing": "true",
                },
            },
            {
                "name": "microbadger",
                "image": "plugins/webhook",
                "settings": {
                    "urls": {
                        "from_secret": "microbadger_url",
                    },
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
            ],
        },
    }]

def gitter(ctx):
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "gitter",
        "clone": {
            "disable": True,
        },
        "steps": [
            {
                "name": "gitter",
                "image": "plugins/gitter",
                "settings": {
                    "webhook": {
                        "from_secret": "gitter_webhook",
                    },
                },
            },
        ],
        "depends_on": [
            "manifest",
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
            ],
            "status": [
                "failure",
            ],
        },
    }]

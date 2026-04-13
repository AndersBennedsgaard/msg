version := `git describe --tags --always`
commit  := `git rev-parse HEAD`
date    := `date -u +'%Y-%m-%dT%H:%M:%SZ'`

repo := "github.com/AndersBennedsgaard/msg/"
versionFlag := repo + "internal/version.Version=" + version
commitFlag := repo + "internal/version.Commit=" + commit
dateFlag := repo + "internal/version.Date=" + date

ldflags := "-X " + versionFlag + " -X " + commitFlag + " -X " + dateFlag

test:
    go test ./...

build:
    go build -ldflags '{{ldflags}}' -o bin/msg main.go

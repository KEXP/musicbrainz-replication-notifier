# musicbrainz-replication-notifier
simple solution to populating a slack channel with nightly mb status.

## Development
* Install Go: [https://go.dev/](https://go.dev/)

## Building
* For Linux: GOOS=linux GOARCH=amd64 go build -o app.go
* For Windows: GOOS=windows GOARCH=amd64 app.go
* For Intel Mac: GOOS=darwin GOARCH=amd64 go build -o app.go
* For Mac-y Mac: GOOS=darwin GOARCH=arm64 go build -o app.go

## Overview
Each night our local musicbrainz will attempt to pull down any new entries to their primary database. The result of this is stored in mirror.log. This simply takes the last 10 lines of that file and sends it to slack. A cron job runs this every morning at 8am PST.

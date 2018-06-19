# drone-github-status

Drone plugin to add Repo Status to a commit.

## Build

Build the binary with the following commands:

```
go test ./...
go build
```

## Docker

Build the docker image with the following commands:

```
docker build -t jmccann/drone-github-status .
```

## Usage

Execute from the working directory:

```
docker run \
  jmccann/drone-github-status --repo-owner jmccann --repo-name drone-github-status \
  --api-key abcd1234 --state "success"
```

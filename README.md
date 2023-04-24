[![Build and Deploy](https://github.com/aorfanos/autommit/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/aorfanos/autommit/actions/workflows/build.yaml)
# autommit

## Description 

Autommit is a simple tool to automatically commit and push changes to a git repository.

It uses the OpenAI API and your local `git` cli to generate appropriate commit messages for your changes.

## Requirements

- OpenAI API key
- `git` executable in your $PATH

## Installation

### Available builds

| OS/Platform  | Architecture | Supported |
|--------------|--------------|-----------|
| Linux        | amd64        | ✅         |
| Linux        | arm64        | ✅         |
| MacOS/Darwin | amd64        | ✅         |
| MacOS/Darwin | arm64        | ✅         |
| Docker       | amd64        | ✅         |
| Docker       | arm64        | TODO       |
| Windows      | amd64        | ❌         |
| Windows      | arm64        | ❌         |

### Install

There are several options to run autommit. Sorting from fastest to more engaged:

1. Fetch binary from this repo's releases page

Get the latest release from https://github.com/aorfanos/autommit/releases/latest/download/autommit-<OS>-<ARCH> and place it in your $PATH.
E.g. for Linux amd64:

```shell
wget https://github.com/aorfanos/autommit/releases/latest/download/autommit-linux-amd64 -O autommit && \
chmod +x autommit && \
sudo mv autommit /usr/local/bin
```

2. Use the Docker image

```shell
docker run -it -e OPENAI_API_KEY=$OPENAI_API_KEY -e GIT_ACC_MAIL="<user-github-mail>" -e GIT_ACC_NAME="<account-fullname)" -v $(pwd):/app:rw autommit
```

3. Build from source

```shell
make build
```
[![Build and Deploy](https://github.com/aorfanos/autommit/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/aorfanos/autommit/actions/workflows/build.yaml)
# autommit

## Description

Autommit is a simple tool to automatically commit and push changes to a git repository.

It uses the OpenAI API and your local `git` cli to generate appropriate commit messages for your changes.

## Requirements

- OpenAI API key
- `git` executable in your $PATH
- If you want to sign/verify your commits, make sure you have a PGP key pair (public and private keys). If you don't have a PGP key pair, you can generate one using [GnuPG](https://gnupg.org/). Then provide the `--pgp-key-path` flag to autommit as follows:

:warning: **Note**: Do **not** use a password for your PGP key. autommit can't handle it (:

```shell
autommit --pgp-key-path /path/to/your/pgp/key
```


## Usage

```shell

❯ autommit --help
Usage of autommit:
  -conventional-commits-type string
        Will add the provided type to the commit message (default "feat")
  -openai-api-key string
        OpenAI API key
  -path string
        Path to the git repository (default ".")
  -pgp-sign
        Will sign the commit with the default PGP key (default true)
  -sign-commits-with-message string
        Will add the provided message to the long commit message (default "Created by autommit 🦄")
  -t string
        Alias of --conventional-commits-type
```

### Demo

- Asciinema [rec](https://goo.com)

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

:warning: **Note**: The Docker image is currently unstable. Prefer the binary builds. Use at your own risk.

```shell
docker run -it -e OPENAI_API_KEY=$OPENAI_API_KEY -e GIT_ACC_MAIL="$(git config user.email)" -e GIT_ACC_NAME="$(git config user.name)" -v $(pwd):/app:rw ghcr.io/aorfanos/autommit:latest
```

3. Build from source

```shell
make build && \
mv autommit-dev /usr/local/bin/autommit
```

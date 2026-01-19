# git-lost-and-found

[![go report card](https://goreportcard.com/badge/github.com/bthuilot/git-lost-and-found/v2)](https://goreportcard.com/report/github.com/bthuilot/git-lost-and-found/v2)
[![GitHub Release](https://img.shields.io/github/v/release/bthuilot/git-lost-and-found)](https://github.com/bthuilot/git-lost-and-found/releases)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


Git scanning tool designed to find dangling commits.

This tool is designed to be used in conjunction with other git scanning tools that leverage `git log` to search diffs.
`git-lost-and-found` is designed to find commits that are not reachable by any branch or tag in the repository, and add named references to them.
This allows other tools to find these commits and their changes, 
since once they are reachable by a named reference they will be included in the output of `git log --all`.
Some tools that can be used in conjunction with `git-lost-and-found` are:

- [gitleaks](https://github.com/gitleaks/gitleaks)
- [trufflehog](https://github.com/trufflesecurity/trufflehog)


## Installing
### Package manager

Install using homebrew
```bash
brew tap bthuilot/tap
brew install bthuilot/tap/git-lost-and-found
```

## Building Locally

The repository can be built using the makefile provided.
Requires Go to be installed on the system.
Optionally a docker image can be built using the makefile.

```bash
# clone the repo
git clone github.com/bthuilot/git-lost-and-found && cd git-lost-and-found

# To build the binary (output in bin/)
make build

# Or to build a docker image (tagged as git-lost-and-found:dev)
DOCKER_IMAGE=git-lost-and-found:dev make docker-build
```

## Running

```bash
# Find dangling commits and don't remove on cleanup
git-lost-and-found find --repo-path . --keep-refs

# Find danling refs, then run gitleaks
# once complete, remove created references
git-lost-and-found find --repo-path . -- gitleaks detect

# clone the linux kernel, find dangling refs,
# run trufflehog, then remove cloneded directory
git-lost-and-found find --repo-url "https://github.com/torvalds/linux" -- trufflehog git file://{} --json

# Use the help menu to see what options are available
git-lost-and-found find --help
```

## CI Script 

A bash script is also provided to enable existing CI
infrastrucre to perform the lost and found lookup for references.
THe only requirements for the script are `sh`, `curl` and `git`.

```bash
# this assumes the cwd is inside a git directory
sh -c "$(curl -fsSL https://git-lf.thuilot.io/ci-scan)"
```

## Example scans


#### Scanning a local git repository with trufflehog (via Docker)

```bash
# git repository cloned to /my/repo/path
docker run -v /my/repo/path:/target \
  ghcr.io/bthuilot/git-lost-and-found:latest find \
  --repo-path /target \
  -- trufflehog git file://. --no-verification
```

### Scanning an existing repo using gitleaks (via CLI)

```bash
# git repository cloned to /my/repo/path
# NOTE: gitleaks will have to be installed on the system
git-lost-and-found find --repo-path "/my/repo/path" \
	-- gitleaks detect .
```

### Clone and scan a repository with Trufflehog (via Docker)
```bash
# NOTE: trufflehog will have to be installed on the system
git-lost-and-found find --repo-url "https://github.com/torvalds/linux" \
	-- trufflehog git file://. --no-verification
```
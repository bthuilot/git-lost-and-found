# git-lost-and-found

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
brew install bthuilot/tap/git-scanner
```

## Building Locally

The repository can be built using the makefile provided.
Requires Go to be installed on the system.
Optionally a docker image can be built using the makefile.

```bash
# clone the repo
git clone github.com/bthuilot/git-scanner && cd git-scanner
# To build the binary
make build
# Or to build a docker image
make build-docker
```

## Running


```bash
# Use the help menu to see what options are available
git-scanner scan --help
```

### Using a docker image
```bash
docker run \
  -v /my/repo/path:/repo \
  ghcr.io/bthuilot/git-scanner:latest scan \
  --repo-path /repo -- trufflehog git file://. --no-verification
```

### Scanning an existing repo using gitleaks

```bash
# NOTE: gitleaks will have to be installed on the system
# git-scanner is not responsible for installing or configuring gitleaks
git-scanner scan --repo-path "/my/repo/path" -- gitleaks detect .
# OR  git-scanner scan --repo-path "/my/repo/path" -- gitleaks detect {}
```

### Clone and scan a repository with Trufflehog
```bash
# NOTE: trufflehog will have to be installed on the system
# git-scanner is not responsible for installing or configuring trufflehog
git-scanner scan --repo-url "https://github.com/torvalds/linux" -- trufflehog git file://. --no-verification
# OR  git-scanner scan --repo-url "https://github.com/torvalds/linux" -- trufflehog git file://{} --no-verification
```
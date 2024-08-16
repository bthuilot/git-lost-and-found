# git-scanner

Git scanning tool designed to also scanning dangling commits, which many other tools miss.
This scanning tool scan run either gitleaks or trufflehog over the full set of commits of a repository

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
  -v /my/repo/path:/repo -v /my/output/path:/output \
  ghcr.io/bthuilot/git-scanner:latest scan \
  --repo-path /repo --scanner "gitleaks" --output /output/results.json
```

### Scanning an existing repo using gitleaks
```bash
git-scanner scan --repo-path "/my/repo/path" --scanner "gitleaks" --output /tmp/results.json --scanner-config ~/gitleaks.toml  
```

### Clone and scan a repository with Trufflehog
```bash
git-scanner scan --repo-url "https://github.com/torvalds/linux" --scanner "trufflehog" --output /tmp/results.json \
  --scanner-args="--no-verification" # Additional args to pass to scanner
```
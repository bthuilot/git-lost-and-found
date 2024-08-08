# git-scanner

Git scanning tool designed to also scanning dangling commits, which many other tools miss.
This scanning tool scan run either gitleaks or trufflehog over the full set of commits of a repository

## Installing
### Package manager
```bash

```

## Building Locally
```bash
make build
# Or to build a docker image
make build-docker
```

## Running
```bash
# Use the help menu to see what options are available
git-scanner scan --help
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
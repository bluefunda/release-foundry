# Homebrew Distribution

release-foundry supports Homebrew installation via a tap repository.

## Tap structure

Create a separate GitHub repository named `homebrew-tap` in your org:

```
<org>/homebrew-tap/
  Formula/
    release-foundry.rb
```

GoReleaser generates and commits the formula automatically on every release.

## Setup

### 1. Create the tap repo

```bash
gh repo create <org>/homebrew-tap --public
```

Initialize with a `Formula/` directory:

```bash
mkdir Formula && touch Formula/.gitkeep
git add Formula/.gitkeep && git commit -m "chore: init tap"
git push
```

### 2. Create a GitHub token for the tap

GoReleaser needs write access to push the generated formula. Create a fine-grained
PAT with **Contents: Read and write** scoped to the `homebrew-tap` repo.

Add it as a repository or org secret named `HOMEBREW_TAP_TOKEN`.

### 3. Configure .goreleaser.yml

The `brews` section in `.goreleaser.yml` is already configured. Verify the
`repository.owner` matches your org:

```yaml
brews:
  - repository:
      owner: '{{ .Env.GITHUB_REPOSITORY_OWNER }}'
      name: homebrew-tap
      token: '{{ .Env.HOMEBREW_TAP_TOKEN }}'
```

### 4. Pass the secret in the release workflow

The `go-binary-release.yml` workflow already accepts `HOMEBREW_TAP_TOKEN`:

```yaml
goreleaser:
  uses: <org>/release-foundry/.github/workflows/go-binary-release.yml@main
  secrets:
    HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

## User installation

Once the tap is published:

```bash
brew tap <org>/tap
brew install release-foundry
```

Or in one command:

```bash
brew install <org>/tap/release-foundry
```

## Upgrading

```bash
brew update && brew upgrade release-foundry
```

## Verifying the installation

```bash
release-foundry version
```

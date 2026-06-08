# Changelog

## [1.6.9](https://github.com/bluefunda/release-foundry/compare/v1.6.8...v1.6.9) (2026-06-08)


### Bug Fixes

* **go-ci:** move checkout before PAT config, add persist-credentials: false ([#56](https://github.com/bluefunda/release-foundry/issues/56)) ([85e8ac1](https://github.com/bluefunda/release-foundry/commit/85e8ac13132f71287bdf7adb6db04620c0d772e9))
* **go-ci:** use http.extraHeader for auth instead of URL-embedded token ([a495be6](https://github.com/bluefunda/release-foundry/commit/a495be645f3f620de6c988e3e465598d16fa7982))
* **release-notes:** don't persist credentials on release-foundry checkout ([0838c28](https://github.com/bluefunda/release-foundry/commit/0838c2842cef55cb92061f9c6020e408100107c4))

## [1.6.8](https://github.com/bluefunda/release-foundry/compare/v1.6.7...v1.6.8) (2026-06-04)


### Bug Fixes

* replace gh release list/view --json with gh api calls ([de960ca](https://github.com/bluefunda/release-foundry/commit/de960ca8821f11727da268c940516dc99ac6ad1b))

## [1.6.7](https://github.com/bluefunda/release-foundry/compare/v1.6.6...v1.6.7) (2026-06-04)


### Bug Fixes

* add platforms input and make QEMU conditional ([5615384](https://github.com/bluefunda/release-foundry/commit/56153841179219aa57bdbfab4bcf58891cf8bc1a))

## [1.6.6](https://github.com/bluefunda/release-foundry/compare/v1.6.5...v1.6.6) (2026-06-03)


### Bug Fixes

* pin binfmt to qemu-v8.1.5 to avoid tonistiigi/binfmt:latest JSON regression ([ad2d683](https://github.com/bluefunda/release-foundry/commit/ad2d683b23674b14b65de68b4ce5f29848d5ea52))
* remove sboms from goreleaser config ([#45](https://github.com/bluefunda/release-foundry/issues/45)) ([47071d4](https://github.com/bluefunda/release-foundry/commit/47071d4e7d3b413d5150ba3042aa44ba8226b9e3))

## [1.6.5](https://github.com/bluefunda/release-foundry/compare/v1.6.4...v1.6.5) (2026-06-02)


### Bug Fixes

* add QEMU for multi-arch builds, remove gitops coupling, update docs ([#42](https://github.com/bluefunda/release-foundry/issues/42)) ([888b26c](https://github.com/bluefunda/release-foundry/commit/888b26c6081f38958db81e4bacfe463b8346d78b))

## [1.6.4](https://github.com/bluefunda/release-foundry/compare/v1.6.3...v1.6.4) (2026-06-01)


### Bug Fixes

* clean up stale /tmp/gitops before cloning in gitops update step ([a295be1](https://github.com/bluefunda/release-foundry/commit/a295be1982c3c0d89efdf68dd94a495b9fe6704e))
* default runner to ubuntu-latest for public/OSS compatibility ([#38](https://github.com/bluefunda/release-foundry/issues/38)) ([9ad23d9](https://github.com/bluefunda/release-foundry/commit/9ad23d9cfa3750811055d4e011ed60b6035416a3))

## [1.6.3](https://github.com/bluefunda/release-foundry/compare/v1.6.2...v1.6.3) (2026-05-31)


### Bug Fixes

* default all reusable workflow runners to self-hosted ([#34](https://github.com/bluefunda/release-foundry/issues/34)) ([e8985bf](https://github.com/bluefunda/release-foundry/commit/e8985bf6bf291d6d2494873272ff6ef23ecc8c80))
* detect and configure rootless Docker socket on self-hosted runners ([#36](https://github.com/bluefunda/release-foundry/issues/36)) ([21b2bdf](https://github.com/bluefunda/release-foundry/commit/21b2bdfb90c19e28617d25ae21f552088c7d42f0))

## [1.6.2](https://github.com/bluefunda/release-foundry/compare/v1.6.1...v1.6.2) (2026-05-31)


### Bug Fixes

* use repository_owner to resolve release-foundry source repo ([#32](https://github.com/bluefunda/release-foundry/issues/32)) ([fe771dc](https://github.com/bluefunda/release-foundry/commit/fe771dccb4daaab09def254b324179530c519981))

## [1.6.1](https://github.com/bluefunda/release-foundry/compare/v1.6.0...v1.6.1) (2026-05-31)


### Bug Fixes

* auto-inject GH_PAT into Docker build-args when provided as secret ([#30](https://github.com/bluefunda/release-foundry/issues/30)) ([a49676a](https://github.com/bluefunda/release-foundry/commit/a49676ad1539f39189b8148546db2ad3e7e88710))

## [1.6.0](https://github.com/bluefunda/release-foundry/compare/v1.5.0...v1.6.0) (2026-05-30)


### Features

* add GoReleaser config, SECURITY.md, CONTRIBUTING.md, Homebrew docs ([8cf89e2](https://github.com/bluefunda/release-foundry/commit/8cf89e243966f65af9781e9b137b210a09dca5c3))
* add renderer interface, topic discovery, and version subcommand ([47507c4](https://github.com/bluefunda/release-foundry/commit/47507c49dc1e973b423aa74706898d693a202c8b))
* generalize workflows — add gitops-repo/gitops-compose-path inputs, remove hardcoded defaults ([13a452b](https://github.com/bluefunda/release-foundry/commit/13a452b3e46063310b61081c7786d0008ec019d9))


### Bug Fixes

* derive release-foundry repo from workflow_ref; fix go build working-directory ([cb4253d](https://github.com/bluefunda/release-foundry/commit/cb4253de9fca72075ed07241170b6e1c8c1a9ae7))
* remove unused id-token permission from go-binary-release — breaks callers ([6a82ac1](https://github.com/bluefunda/release-foundry/commit/6a82ac1fc95a9e6db61c4f6cf057ba38a242df18))

## [1.5.0](https://github.com/bluefunda/release-foundry/compare/v1.4.0...v1.5.0) (2026-05-26)


### Features

* route remaining CI/CD jobs to self-hosted runners ([#27](https://github.com/bluefunda/release-foundry/issues/27)) ([9e9779c](https://github.com/bluefunda/release-foundry/commit/9e9779c4523a8d8ac25f70046fda29f03de408fd))


### Bug Fixes

* resolve golangci-lint errors surfaced by new CI ([#25](https://github.com/bluefunda/release-foundry/issues/25)) ([6a72f9f](https://github.com/bluefunda/release-foundry/commit/6a72f9fb974abc09cf3f26d80cfa3d7b4ea2cc59))

## [1.4.0](https://github.com/bluefunda/release-foundry/compare/v1.3.1...v1.4.0) (2026-05-26)


### Features

* route go-ci jobs to self-hosted runner by default ([#23](https://github.com/bluefunda/release-foundry/issues/23)) ([0d1678c](https://github.com/bluefunda/release-foundry/commit/0d1678c77f7242377e70947319d0c75b73949b52))

## [1.3.1](https://github.com/bluefunda/release-foundry/compare/v1.3.0...v1.3.1) (2026-05-23)


### Bug Fixes

* correct Komodo deploy URL to gitops.bluefunda.com ([2d678f6](https://github.com/bluefunda/release-foundry/commit/2d678f6f410ee1e4288721da12c5cc40ba072e36))
* replace Komodo direct API call with GitOps-based deploy ([#20](https://github.com/bluefunda/release-foundry/issues/20)) ([9304cb8](https://github.com/bluefunda/release-foundry/commit/9304cb8414e639f5f2820d3a67f5daec5a641a28))
* show Komodo API response in deploy step for debugging ([829614a](https://github.com/bluefunda/release-foundry/commit/829614a4824e16e757fffd1199b13d0f416bf8a2))

## [1.3.0](https://github.com/bluefunda/release-foundry/compare/v1.2.0...v1.3.0) (2026-05-14)


### Features

* add komodo deploy step after image push ([#16](https://github.com/bluefunda/release-foundry/issues/16)) ([471e721](https://github.com/bluefunda/release-foundry/commit/471e721fd99e576309c3dd178481a42ddc7e44b2))

## [1.2.0](https://github.com/bluefunda/release-foundry/compare/v1.1.0...v1.2.0) (2026-05-14)


### Features

* add github-release-notes reusable workflow and -since flag ([2d7dfc2](https://github.com/bluefunda/release-foundry/commit/2d7dfc27fe2ad9b86cdbabc4025a3cac830c1fa2))

## [1.1.0](https://github.com/bluefunda/release-foundry/compare/v1.0.0...v1.1.0) (2026-05-13)


### Features

* add macOS notarization support to go-binary-release workflow ([#14](https://github.com/bluefunda/release-foundry/issues/14)) ([c888bd1](https://github.com/bluefunda/release-foundry/commit/c888bd16d6c80ae50effe725aaa1e036b2217bc9))
* **go-ci:** add goprivate input, remove hardcoded GOPRIVATE env ([#11](https://github.com/bluefunda/release-foundry/issues/11)) ([a7761ff](https://github.com/bluefunda/release-foundry/commit/a7761ff9badbd78b09af8e45f7ed1a3ad157f485))


### Bug Fixes

* revert Configure Git step to shell if-block ([#13](https://github.com/bluefunda/release-foundry/issues/13)) ([781f230](https://github.com/bluefunda/release-foundry/commit/781f230694054fc6bf94085d45db6f93a11c9dc6))

## 1.0.0 (2026-04-27)


### Features

* add github-release renderer and Phase 1 release pipeline ([#9](https://github.com/bluefunda/release-foundry/issues/9)) ([75418f4](https://github.com/bluefunda/release-foundry/commit/75418f47d4778083f91f23a3d3bd95ba6873b8c5))
* add go-binary-release reusable workflow ([#1](https://github.com/bluefunda/release-foundry/issues/1)) ([9952371](https://github.com/bluefunda/release-foundry/commit/9952371271db8650ad4052597b2478713fdac9ac))
* add multi-repo batch mode with edition support ([#4](https://github.com/bluefunda/release-foundry/issues/4)) ([fdcbf8e](https://github.com/bluefunda/release-foundry/commit/fdcbf8e090755879729cb10160772b2a2c94b407))
* add reusable Go CI and Docker deploy workflows ([6f484b9](https://github.com/bluefunda/release-foundry/commit/6f484b9536bcd8db54dbf8eb575c4e06af23f420))
* add workflow_dispatch to docker-deploy for manual image builds ([123d1ab](https://github.com/bluefunda/release-foundry/commit/123d1ab0074b6a78f2d093ee95322f8e360d161e))
* initial release-foundry service ([f576624](https://github.com/bluefunda/release-foundry/commit/f57662419a899ed4ba3fcea85577ccc4a05b1098))
* standardize release pipeline with release-please ([#2](https://github.com/bluefunda/release-foundry/issues/2)) ([e09e31f](https://github.com/bluefunda/release-foundry/commit/e09e31fb52fd83e84564346ccc01b96568547760))


### Bug Fixes

* bump golangci-lint default to v2.11.4 (Go 1.24+ support) ([b1c82c9](https://github.com/bluefunda/release-foundry/commit/b1c82c92321091bcf3b84a455ca888cf50f27272))
* checkout specific tag ref in go-binary-release to support manual dispatch ([dd9e2ec](https://github.com/bluefunda/release-foundry/commit/dd9e2ec53a1c9ed1fa9a510b5d25eb0025eec45b))
* gate GH_PAT git config via shell, not step-level if ([#7](https://github.com/bluefunda/release-foundry/issues/7)) ([01ee057](https://github.com/bluefunda/release-foundry/commit/01ee0579ffa4067f649a1ad6b1896a3618b55a8b))
* gate private-module git config via shell, not step-level if ([#8](https://github.com/bluefunda/release-foundry/issues/8)) ([a37a6dc](https://github.com/bluefunda/release-foundry/commit/a37a6dc01e7ba410086a4e1765fd67fba46ffd3b))
* make GH_PAT optional in go-binary-release workflow ([6aa9d65](https://github.com/bluefunda/release-foundry/commit/6aa9d659065c815184e639e2464c65924e2cc787))
* rename HOMEBREW_TAP_GITHUB_TOKEN to HOMEBREW_TAP_TOKEN in goreleaser env ([138e326](https://github.com/bluefunda/release-foundry/commit/138e326245a9ab1ab7589376084e68bd355cd3f5))
* upgrade release-please-action to v4.1.3, use Node.js 24 ([43858f2](https://github.com/bluefunda/release-foundry/commit/43858f26e389a1ccebc0d81d29705bc54658d543))
* use GH_PAT in release-please to trigger downstream CI ([#5](https://github.com/bluefunda/release-foundry/issues/5)) ([93d5ea8](https://github.com/bluefunda/release-foundry/commit/93d5ea8d77fd14b5150e98d4e6786c506cf64b90))

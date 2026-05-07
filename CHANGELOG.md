# Changelog

## [1.1.0](https://github.com/bluefunda/release-foundry/compare/v1.0.0...v1.1.0) (2026-05-07)


### Features

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

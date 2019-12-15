# golang-starter
[![Build Status](https://travis-ci.com/chenpengfei/golang-starter.svg)](https://travis-ci.com/chenpengfei/golang-starter)
[![Coverage Status](https://coveralls.io/repos/github/chenpengfei/golang-starter/badge.svg)](https://coveralls.io/github/chenpengfei/golang-starter)

Golang Starter，用于创建符合 Semantic Release 规范的初始仓库

## Commit
### Pre Commit(Staged Code)
1. 格式化工具 -> gofmt
2. 语义化消息 -> [commitizen](https://github.com/commitizen/cz-cli)
### Commit Message
1. 校验消息 -> [commitlint])(https://github.com/conventional-changelog/commitlint)
2. Git hooks -> [husky](https://github.com/typicode/husky)
### Pre Push
1. Lint -> golint + golangci-lint
2. Unit Test -> go test


## Configuration
1. `makefile` 组合命令
2. `.editorconfig` 文档格式配置

## Link
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)
[Semantic-Release](./docs/Semantic-Release.key)
[SemVer](https://semver.org/lang/zh-CN/)
[Golang SemVer](https://golang.org/src/cmd/go/internal/semver/semver.go)
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

## FAQ
`make init` 提示错误 `rollbackFailedOptional: verb npm-session 8c37f8a1a41ff065`，执行以下步骤解决
```
> npm config rm proxy
> npm config rm https-proxy
> npm config set registry http://registry.npmjs.org
```


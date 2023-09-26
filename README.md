<!-- <p align="center"><img src="" alt="" height="100px"></p> -->

<div align="center">
  <a href="https://codecov.io/gh/lsytj0413/golang-project-template" > 
    <img src="https://codecov.io/gh/lsytj0413/golang-project-template/branch/master/graph/badge.svg?token=XM1YHY2D3R"/> 
  </a>
  <a href="https://github.com/lsytj0413/golang-project-template/actions">
    <img src="https://github.com/lsytj0413/golang-project-template/workflows/Unit%20tests/badge.svg" alt="Actions Status">
  </a>
  <a href="https://goreportcard.com/report/github.com/lsytj0413/golang-project-template">
    <img src="https://goreportcard.com/badge/github.com/lsytj0413/golang-project-template?style=flat-square" alt="Go Report Card">
  </a>
</div>

# golang-project-template

## What is this？

This is an template for `golang` project, it has the following features:

+ An `Makefile` style for build、test，support docker
+ Support addlicense headers
+ Buildin grpc、grpc-gateway with example
+ Github template and action for PR、ISSUE and workflow
+ Go 1.21，with golangci-lint、mockgen and so on, and unittest is enabled
+ Codecov enabled

## How to use？

+ First, create an repository with this template, and clone it
+ Second, run `./hack/rename {xxx}` to new project name
+ Third, follow [codecov](https://about.codecov.io/) to configration your `CODECOV_TOKEN`

Now, enjoy coding!
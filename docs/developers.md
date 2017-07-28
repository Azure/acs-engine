# Developers Guide

This guide explains how to set up your environment for developing on
acs-engine.

## Prerequisites

- Go 1.6.0 or later
- Glide 0.12.0 or later
- kubectl 1.5 or later
- An Azure account (needed for deploying VMs and Azure infrastructure)
- Git

## Contribution Guidelines

We welcome contributions. This project has set up some guidelines in
order to ensure that (a) code quality remains high, (b) the project
remains consistent, and (c) contributions follow the open source legal
requirements. Our intent is not to burden contributors, but to build
elegant and high-quality open source code so that our users will benefit.

Make sure you have read and understood the main CONTRIBUTING guide:

https://github.com/Azure/acs-engine/blob/master/CONTRIBUTING.md

### Structure of the Code

The code for the acs-engine project is organized as follows:

- The individual programs are located in `cmd/`. Code inside of `cmd/`
  is not designed for library re-use.
- Shared libraries are stored in `pkg/`.
- The `tests/` directory contains a number of utility scripts. Most of these
  are used by the CI/CD pipeline.
- The `docs/` folder is used for documentation and examples.

Go dependencies are managed with
[Glide](https://github.com/Masterminds/glide) and stored in the
`vendor/` directory.

### Git Conventions

We use Git for our version control system. The `master` branch is the
home of the current development candidate. Releases are tagged.

We accept changes to the code via GitHub Pull Requests (PRs). One
workflow for doing this is as follows:

1. Use `go get` to clone the acs-engine repository: `go get github.com/Azure/acs-engine`
2. Fork that repository into your GitHub account
3. Add your repository as a remote for `$GOPATH/github.com/Azure/acs-engine`
4. Create a new working branch (`git checkout -b feat/my-feature`) and
   do your work on that branch.
5. When you are ready for us to review, push your branch to GitHub, and
   then open a new pull request with us.

### Go Conventions

We follow the Go coding style standards very closely. Typically, running
`go fmt` will make your code beautiful for you.

We also typically follow the conventions recommended by `go lint` and
`gometalinter`. Run `make test-style` to test the style conformance.

Read more:

- Effective Go [introduces formatting](https://golang.org/doc/effective_go.html#formatting).
- The Go Wiki has a great article on [formatting](https://github.com/golang/go/wiki/CodeReviewComments).

# my-cargonaut/cargonaut

> Educational ride and transport sharing app for an university course. - by **[Lukas Malkmus], [Philipp Alexander Händler] & [Robert Feuerhack]**

[![Build Status][build_badge]][build]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![GoDoc][docs_badge]][docs]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]
[![License Status][license_status_badge]][license_status]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

*Cargonaut* is an dducational ride and transport sharing app for an university
course.

## Usage

### Installation

This project uses native [go mod] support for vendoring and requires a working
`go` toolchain installation.

Binary releases are currently not available.

#### Install using `go get`

```bash
GO111MODULE=on go install github.com/my-cargonaut/cargonaut/cmd/...
```

#### Install from source

```bash
git clone https://github.com/my-cargonaut/cargonaut.git
cd cargonaut
make # Build production binaries
make install # Build and install binaries into $GOPATH
```

#### Validate installation

The installation can be validated by running `cargonaut version` in the terminal.

### Configuration

The application can be configured by either the environment or command line
flags. It provides a basic *help flag* `-help` which prints out application
and configuration help. See [Using the application](#using-the-application).

The application can be configured by command line flags:

```bash
cargonaut serve -listen-address=:8080
```

Flags are mapped to environment variables by replacing dashes `-` with an 
underscore `_`. However, they are prefixed by the application name:

```bash
export CARGONAUT_LISTEN_ADDRESS=:8080
```

Configuration priority from lowest to highest is like presented above:
Command line option (flag), environment.

### Using the application

```bash
cargonaut [flags] [commands]
```

Help on flags and commands:

```bash
cargonaut --help
```

## Contributing

Feel free to submit PRs or to fill issues. Every kind of help is appreciated.

Before committing, `make` should pass without any issues.

## License

© Lukas Malkmus, Philipp Alexander Händler, Robert Feuerhack 2020

Distributed under MIT License (`The MIT License`).

<!-- Links -->

[Lukas Malkmus]: https://github.com/lukasmalkmus
[Philipp Alexander Händler]: https://github.com/philippalexanderhaendler
[Robert Feuerhack]: https://github.com/RFeuerhack
[go mod]: https://golang.org/cmd/go/#hdr-Module_maintenance

<!-- Badges -->

[build]: https://travis-ci.com/my-cargonaut/cargonaut
[build_badge]: https://img.shields.io/travis/com/my-cargonaut/cargonaut.svg?style=flat-square
[coverage]: https://codecov.io/gh/my-cargonaut/cargonaut
[coverage_badge]: https://img.shields.io/codecov/c/github/my-cargonaut/cargonaut.svg?style=flat-square
[report]: https://goreportcard.com/report/github.com/my-cargonaut/cargonaut
[report_badge]: https://goreportcard.com/badge/github.com/my-cargonaut/cargonaut?style=flat-square
[docs]: https://godoc.org/github.com/my-cargonaut/cargonaut
[docs_badge]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square
[release]: https://github.com/my-cargonaut/cargonaut/releases
[release_badge]: https://img.shields.io/github/release/my-cargonaut/cargonaut.svg?style=flat-square
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/my-cargonaut/cargonaut.svg?color=blue&style=flat-square
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Fmy-cargonaut%2Fcargonaut?ref=badge_shield
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmy-cargonaut%2Fcargonaut.svg
[license_status_large]: https://app.fossa.com/projects/git%2Bgithub.com%2Fmy-cargonaut%2Fcargonaut?ref=badge_large
[license_status_large_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmy-cargonaut%2Fcargonaut.svg?type=large

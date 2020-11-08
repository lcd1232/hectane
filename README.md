## Hectane

[![Go Report Card](https://goreportcard.com/badge/github.com/lcd1232/hectane)](https://goreportcard.com/report/github.com/lcd1232/hectane)
[![Build Status](https://img.shields.io/github/workflow/status/lcd1232/hectane/test)](https://github.com/lcd1232/hectane/actions)
[![Coverage](https://img.shields.io/codecov/c/github/lcd1232/hectane)](https://codecov.io/gh/lcd1232/hectane)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/github.com/lcd1232/hectane)
[![Releases](https://img.shields.io/github/v/tag/lcd1232/hectane.svg)](https://github.com/lcd1232/hectane/releases)
[![LICENSE](https://img.shields.io/github/license/lcd1232/hectane.svg)](https://github.com/lcd1232/hectane/blob/master/LICENSE)

Hectane is both a Go package providing an SMTP queue for sending emails and a standalone application that exposes this functionality via an HTTP API.

### Features

- Ability to attach files to emails
- Support for TLS encryption and HTTP basic auth
- Mail queue that efficiently delivers emails to hosts
- Emails in the queue are stored on disk until delivery
- MX records for the destination host are tried in order of priority
- Run the application as a service on Windows

### Documentation

Documentation for Hectane can be found below:

- [Using Hectane in a Go application](https://github.com/hectane/hectane/wiki/Hectane%20Package)
- [Using Hectane in another language or on a server](https://github.com/hectane/hectane/wiki/Hectane%20Daemon)

### Installation

In addition to the [files on the releases page](https://github.com/hectane/hectane/releases), Hectane can be installed from any of the sources below:

- PPA: [stable](https://launchpad.net/~hectane/+archive/ubuntu/hectane) | [daily](https://launchpad.net/~hectane/+archive/ubuntu/hectane-dev)
- [Juju charm store](https://jujucharms.com/hectane/)
- [Docker Hub](https://hub.docker.com/r/hectane/hectane/)

# mbgo

[![GoDoc](https://godoc.org/github.com/ogbofjnr/mbgo?status.svg)](https://godoc.org/github.com/ogbofjnr/mbgo) [![Build Status](https://travis-ci.org/senseyeio/mbgo.svg?branch=master)](https://travis-ci.org/senseyeio/mbgo) [![Go Report Card](https://goreportcard.com/badge/github.com/ogbofjnr/mbgo)](https://goreportcard.com/report/github.com/ogbofjnr/mbgo)

A mountebank API client for the Go programming language.

## Installation

```sh
go get -u github.com/ogbofjnr/mbgo@latest
```

## Testing

This package includes both unit and integration tests. Use the `unit` and `integration` targets in the Makefile to run them, respectively:

```sh
make unit
make integration
```

The integration tests expect Docker to be available on the host, using it to run a local mountebank container at 
`localhost:2525`, with the additional ports 8080-8081 exposed for test imposters. Currently tested against a mountebank 
v2.1.2 instance using the [andyrbell/mountebank](https://hub.docker.com/r/andyrbell/mountebank) image on DockerHub.

## Contributing

* Fork the repository.
* Code your changes.
* If applicable, add tests and/or documentation.
* Please ensure all unit and integration tests are passing, and that all code passes `make lint`.
* Raise a new pull request with a short description of your changes.
* Use the following convention for branch naming: `<username>/<description-with-dashes>`. For instance, `smotes/add-smtp-imposters`.

# Metal Makefile

As you can probably guess, the development and CI/CD workflow for Metal is
driven with Make. This document will contain sections related to build targets
and recipes. To have a look at this documentation in your local browser,
execute:

```make```

## Development Environments

For each service and server within metal there is a corresponding development
environment driven by docker-compose files in the docker directory. A
development environment is defined by a yml file within this directory without
the file extension. The current list of development environments:

* rexecd-mysql-server

There are also actions associated with each development environment. The
following actions are:

* start
* status
* stop
* restart

For each development environment and action there exists a phony make target
for orchestrating your development environment:


* rexecd-mysql-server-start
* rexecd-mysql-server-status
* rexecd-mysql-server-stop
* rexecd-mysql-server-restart

## Testing

* unit: executes all unit tests
* smoke: executes all smoke tests

## Protobuf

For each service managed via gRPC there is a corresponding directory under
proto containing protobuf files defining services and messages for the
corresponding service.

* proto: rebuild all golang packages from the corresponding protobuf files

## Dependencies

All dependencies are managed with Go modules.

* get: update all golang dependencies. You may pass the -e flag to make and set
       the PACKAGE variable to only get one package
* update: same as get, but updates

## Installation

Local installation is handled via the go binary

* install: installs metal into GOBIN

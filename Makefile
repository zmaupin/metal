
.PHONY: clean install proto memory_start memory_status memory_stop build_info smoke app

SHELL := $(shell which bash)
VERSION := $(shell cat VERSION)
PROTO_DIR := proto
# Each proto service will have a directory under proto
PROTO_SERVICES := $(shell find proto -type d -mindepth 1 -exec basename {} +)

# get name of development environments based on docker-compose files found in the docker dir
DEV_ENVS := $(shell find docker -name "*.yml" | xargs basename | sed 's/\.yml//')
DEV_ENV_ACTIONS := start status stop restart

# The first task should always show information on how to use make. This will
# build and open up the documentation in a browser
info: metalmakedocs
	@open docs/make.html

metalmakedocs: build_metalmakedocs
	@metalmakedocs

build_metalmakedocs:
	@go install github.com/metal-go/metal/tools/metalmakedocs

# List of packages to test. Use the elipses to decend into sub packages
PACKAGES := ./cmd... ./rexecd... ./util...

app:
	@cd app && find dist -type f ! -name "*index.html" -exec rm -rf {} + && npx webpack;

# dev_start starts a given development environment
#
# 1: development environment name
define dev_start
	docker-compose --file docker/${1}.yml up --remove-orphans --detach --force-recreate --build
endef

# dev_status retrieves the status of a given development environment
#
# 1: development environment name
define dev_status
	docker-compose --file docker/${1}.yml ps
endef

# dev_stop stops a given development environment
#
# 1: development environment name
define dev_stop
	docker-compose --file docker/${1}.yml stop
endef

# dev_restart forces recreation of a development environment
#
# 1: development environment name
define dev_restart
	docker-compose --file docker/${1}.yml build
	docker-compose --file docker/${1}.yml up --detach --force-recreate
endef

# dev_env_template will be evaluated to create dynamically generated targets
#
# 1: environment
# 2: action
define dev_env_template
${1}-${2}:
	$(call dev_${2},${1})
endef

# For each dev env and action, create a phony target that executes the
# corresponding func.
$(foreach env,${DEV_ENVS},$(foreach action,${DEV_ENV_ACTIONS},$(eval $(call dev_env_template,${env},${action}))))

clean:
	@go clean

proto_setup:
	@go install github.com/golang/protobuf/protoc-gen-go
	@which protoc > /dev/null

# Build all go packages for proto-based services
proto: proto_setup
	@for svc in ${PROTO_SERVICES}; do \
		protoc \
			--proto_path=${PROTO_DIR}/${svc}                           \
			--go_out=plugins=grpc:${PROTO_DIR}/${svc}                  \
			$(shell find ${PROTO_DIR}/${svc} -type f -name "*.proto"); \
	done

install: clean app
	@bash scripts/install

unit:
	@export GOCACHE=off; go test -v ${PACKAGES}

smoke:
	@export GOCACHE=off; go test -v ./test/smoke/rexecdtest...

# PACKAGE may be set with -e
get:
	go get ${PACKAGE}

# PACKAGE may be set with -e
upgrade:
	go get -u ${PACKAGE}

verify:
	@go mod verify

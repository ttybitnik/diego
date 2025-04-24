# MKDEV 0.5.0 (x-release-please-version)
# See <https://github.com/ttybitnik/mkdev> for more information.

PROJECT_NAME = diego
CONTAINER_ENGINE = podman

RUN_BIND_SOCKET = false
EXEC_SHELL_CMD = /bin/bash

__USER = $(or $(USER),$(shell whoami))
__SOCKET = /run/user/$(shell id -u)/podman/podman.sock

# Host targets/commands
.PHONY: dev start open stop clean serestore

dev:
	$(info Building development container image...)

	$(CONTAINER_ENGINE) build \
	--build-arg USERNAME=$(__USER) \
	-f .mkdev/Containerfile \
	-t localhost/mkdev/$(PROJECT_NAME) \
	.

start:
	$(info Starting development container...)

	$(CONTAINER_ENGINE) run -it -d --replace \
	$(if $(filter podman,$(CONTAINER_ENGINE)),--userns=keep-id) \
	--name mkdev-$(PROJECT_NAME) \
	--volume .:/home/$(__USER)/workspace:Z \
	$(if $(filter true,$(RUN_BIND_SOCKET)),--volume $(__SOCKET):$(__SOCKET)) \
	$(if $(filter true,$(RUN_BIND_SOCKET)),--env CONTAINER_HOST=unix://$(__SOCKET)) \
	localhost/mkdev/$(PROJECT_NAME):latest

	@# $(CONTAINER_ENGINE) compose .mkdev/compose.yaml up -d

open:
	$(info Opening development container...)

	$(CONTAINER_ENGINE) exec -it mkdev-$(PROJECT_NAME) $(EXEC_SHELL_CMD)

stop:
	$(info Stopping development container...)

	$(CONTAINER_ENGINE) stop mkdev-$(PROJECT_NAME)

	@# $(CONTAINER_ENGINE) compose .mkdev/compose.yaml down

clean: distclean
	$(info Removing development container and image...)

	-$(CONTAINER_ENGINE) rm mkdev-$(PROJECT_NAME)
	-$(CONTAINER_ENGINE) image rm localhost/mkdev/$(PROJECT_NAME)

	@# $(CONTAINER_ENGINE) image prune

serestore:
	$(info Restoring project SELinux context and permissions...)

	chcon -Rv unconfined_u:object_r:user_home_t:s0 .
	# find . -type d -exec chmod 700 {} \;
	# find . -type f -exec chmod 600 {} \;

# Container targets/commands
.PHONY: lint test build run deploy debug distclean

lint:
	$(info Running linters...)

	golangci-lint run

test: lint
	$(info Running tests...)

	go test -cover -race ./internal/adapters/left/cli

build: test
	$(info Building...)

	go build -o $(PROJECT_NAME) main.go

run: build
	$(info Running...)

	./$(PROJECT_NAME)

deploy: build
	$(info Deploying...)

	goreleaser build --snapshot --clean

debug: test
	$(info Debugging tasks...)

	go build -gcflags=all="-N -l" -o $(PROJECT_NAME) main.go

distclean:
	$(info Cleaning artifacts...)

	rm -rf ./$(PROJECT_NAME) ./dist/


update:
	go get -u ./...
	go mod tidy
	$(MAKE) run

golden: lint
	go test -cover ./internal/adapters/left/cli -update
	$(MAKE) run

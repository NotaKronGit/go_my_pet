.PHONY: all build coverage

COMMANDS = $(foreach DIR, $(dir $(wildcard app/services/*/.)), $(shell basename $(DIR)) )
PKG_LIST = $(shell go list ./... | grep -v /vendor/ | grep -v autotests) # sync with coverage.sh
BUILD_DATE = $(shell LC_TIME=C date)


dep:
	go mod tidy -v
	go mod vendor


.PHONY: lint
lint: ; golangci-lint run --config .golangci.yml

build-mac: dep
	@mkdir -p bin
	for COMMAND in ${COMMANDS} ; do \
		GODEBUG=asyncpreemptoff=1 GO111MODULE=on go build -mod=vendor -v -o bin/$${COMMAND}_mac -ldflags '-X "${PKG_VERSION}.Version=${GIT_VERSION}" -X "${PKG_VERSION}.VersionMajor=${GIT_VERSION_MAJOR}" -X "${PKG_VERSION}.VersionMinor=${GIT_VERSION_MINOR}" -X "${PKG_VERSION}.VersionPatch=${GIT_VERSION_PATCH}" -X "${PKG_VERSION}.SHA=${GIT_SHA}" -X "${PKG_VERSION}.BuildDate=${BUILD_DATE}"' app/services/$${COMMAND}/*.go \
		|| exit 1 ; \
	done
build-linux: dep
	@rm -fr bin
	@mkdir -p bin
	for COMMAND in ${COMMANDS} ; do \
		GODEBUG=asyncpreemptoff=1 GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -v -o bin/$${COMMAND}_linux -ldflags '-X "${PKG_VERSION}.Version=${GIT_VERSION}" -X "${PKG_VERSION}.VersionMajor=${GIT_VERSION_MAJOR}" -X "${PKG_VERSION}.VersionMinor=${GIT_VERSION_MINOR}" -X "${PKG_VERSION}.VersionPatch=${GIT_VERSION_PATCH}" -X "${PKG_VERSION}.SHA=${GIT_SHA}" -X "${PKG_VERSION}.BuildDate=${BUILD_DATE}"' app/services/$${COMMAND}/*.go \
		|| exit 1 ; \
	done


build-linux-image:
	DOCKER_BUILDKIT=1 docker build --platform=linux/amd64 -t go_my_pet_user -f ci/Dockerfile .


restart: pull down up
up: ; docker-compose -f docker-compose.yml up -d
run: down up
down: ; docker-compose -f docker-compose.yml down
ps: ; docker-compose ps
pull: ; docker-compose -f docker-compose.yml pull

up-all: ; docker-compose -f docker-compose.yml up -d
down-all: ;  docker-compose -f docker-compose.yml down
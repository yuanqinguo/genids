export GO111MODULE=on
export VERSION=1.0.0
export ENV=prod
export PROJECT=genids
export CGO_ENABLED=0
export PROJECT_COMMIT_SHA=`git describe --always`

OBJTYPE=api
TOPDIR=$(shell pwd)
BUILD_TIME=`date +%Y%m%d%H%M%S`
OBJTAR=$(PROJECT).tar.gz

SOURCE_MAIN_FILE=main.go
SOURCE_BINARY_DIR=$(OBJTYPE)-$(PROJECT)
SOURCE_BINARY_FILE=$(SOURCE_BINARY_DIR)/$(OBJTYPE)-$(PROJECT)-$(BUILD_TIME)-$(PROJECT_COMMIT_SHA)

BUILD_FLAG=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X 'main.goVersion=`go version`' -X main.gitCommitid=$(PROJECT_COMMIT_SHA)"

all: domod build pack
	@echo "ALL DONE"
	@echo "Program:       "  $(PROJECT)
	@echo "Version:       "  $(VERSION)
	@echo "Env:           "  $(ENV)
	@echo "Commitid:      "  $(PROJECT_COMMIT_SHA)

build:
	@echo "start go build...."$(TOPDIR)
	@rm -rf $(SOURCE_BINARY_DIR)/*
	@go build $(BUILD_FLAG) -o $(SOURCE_BINARY_FILE) $(SOURCE_MAIN_FILE)

domod:
	@echo "start go mod tidy...."$(TOPDIR)
	@go mod tidy
	@echo "go mod graph...."$(PROJECT)
	@go mod graph

pack:
	@echo "packing....tar czvf $(OBJTAR) $(SOURCE_BINARY_DIR)"
	@tar czvf $(OBJTAR) $(SOURCE_BINARY_DIR)

# Douyu Golang Application Standard Makefile

SHELL:=/bin/bash
BASE_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_PATH:=$(BASE_PATH)/build
TITLE:=$(shell basename $(BASE_PATH))
VCS_INFO:=$(shell $(BUILD_PATH)/script/shell/vcs.sh)
BUILD_TIME:=$(shell date +%Y-%m-%d--%T)
APP_PKG:=$(shell $(BUILD_PATH)/script/shell/apppkg.sh)
JUPITER:=$(APP_PKG)/vendor/github.com/labstack/echo/v4/application
LDFLAGS:="-X $(JUPITER).vcsInfo=$(VCS_INFO) -X $(JUPITER).buildTime=$(BUILD_TIME) -X $(JUPITER).name=$(APP_NAME) -X $(JUPITER).id=$(APP_ID)"

all:print fmt buildAgent

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASE_PATH:$(BASE_PATH)
	@echo BUILD_PATH:$(BUILD_PATH)
	@echo TITLE:$(TITLE)
	@echo VCS_INFO:$(VCS_INFO)
	@echo BUILD_TIME:$(BUILD_TIME)
	@echo JUPITER:$(JUPITER)
	@echo APP_NAME:$(APP_NAME)
	@echo -e "\n"

fmt:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making fmt<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	go fmt $(BASE_PATH)/pkg/...
	@echo -e "\n"

lint:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making lint<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
ifndef REVIVE
	go get -u github.com/mgechev/revive
endif
	@revive -formatter stylish pkg/...
	@echo -e "\n"

errcheck:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making errcheck<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
ifndef ERRCHCEK
	go get -u github.com/kisielk/errcheck
endif
	@errcheck pkg/...
	@echo -e "\n"

test:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making test<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo testPath ${BAST_PATH}
	go test -v .${BAST_PATH}/...

buildAgent:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	chmod +x $(BUILD_PATH)/script/shell/*.sh
	$(BUILD_PATH)/script/shell/build.sh $(LDFLAGS)
	@echo -e "\n"

license: ## Add license header for all code files
	@find . -name \*.go -exec sh -c "if ! grep -q 'LICENSE' '{}'; then mv '{}' tmp && cp doc/LICENSEHEADER.txt '{}' && cat tmp >> '{}' && rm tmp; fi" \;


run:
	go run cmd/agent/main.go --config=config/config.toml

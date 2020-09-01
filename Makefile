# Douyu Golang Application Standard Makefile

SHELL:=/bin/bash
BASE_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SCRIPT_PATH:=$(BASE_PATH)/script
APP_NAME:=$(shell basename $(BASE_PATH))
COMPILE_OUT:=$(BASE_PATH)/release
APP_VERSION:=0.4.0


all:print fmt buildAgent

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASE_PATH:$(BASE_PATH)
	@echo SCRIPT_PATH:$(SCRIPT_PATH)
	@echo APP_NAME:$(APP_NAME)
	@echo COMPILE_OUT:$(COMPILE_OUT)
	@echo APP_VERSION:$(APP_VERSION)
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



license: ## Add license header for all code files
	@find . -name \*.go -exec sh -c "if ! grep -q 'LICENSE' '{}'; then mv '{}' tmp && cp doc/LICENSEHEADER.txt '{}' && cat tmp >> '{}' && rm tmp; fi" \;

run:export REGION_CODE=wuhan_region
run:export REGION_NAME=在Agent环境变量修改该参数
run:export ZONE_NAME=光谷
run:export ZONE_CODE=guanggu
run:export ENV=dev
run:
	go run cmd/juno-agent/main.go --config=config/config.toml


build_all:build_agent build_data tar


build_agent:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build juno agent<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@chmod +x $(SCRIPT_PATH)/build/*.sh
	@cd cmd/juno-agent && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT) $(APP_VERSION)
	@echo -e "\n"

build_data:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build juno agent data<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@chmod +x $(SCRIPT_PATH)/build/*.sh
	@$(SCRIPT_PATH)/build/build_data.sh $(APP_NAME) $(APP_VERSION) $(BASE_PATH) $(COMPILE_OUT)/$(APP_VERSION)
	@echo -e "\n"
tar:
	@cd $(BASE_PATH)/release && tar zcvf juno-agent_$(APP_VERSION).tar.gz $(APP_VERSION)

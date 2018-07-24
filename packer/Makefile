SHELL += -eu

BLUE	:= \033[0;34m
GREEN	:= \033[0;32m
RED   := \033[0;31m
NC    := \033[0m

# ∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨∨

export AZURE_LOCATION       		 			?= West US 2
export AZURE_RESOURCE_GROUP_NAME	    ?= aksimages
export AZURE_VM_SIZE 									?= Standard_D2_v2

build:
	@packer build -var-file=settings.json vhd-image.json

init:
	@./scripts/init-variables

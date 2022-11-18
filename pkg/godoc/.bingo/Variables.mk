# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.4.3. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.4.3
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.4.3"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.4.3 "github.com/bwplotka/bingo"

GODOC := $(GOBIN)/godoc-v0.1.2
$(GODOC): $(BINGO_DIR)/godoc.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/godoc-v0.1.2"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=godoc.mod -o=$(GOBIN)/godoc-v0.1.2 "golang.org/x/tools/cmd/godoc"


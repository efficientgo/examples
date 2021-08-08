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
BINGO := $(GOBIN)/bingo-v0.3.1
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.3.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.3.1 "github.com/bwplotka/bingo"

BUF := $(GOBIN)/buf-v0.48.2
$(BUF): $(BINGO_DIR)/buf.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buf-v0.48.2"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=buf.mod -o=$(GOBIN)/buf-v0.48.2 "github.com/bufbuild/buf/cmd/buf"

PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc-v1.1.0
$(PROTOC_GEN_GO_GRPC): $(BINGO_DIR)/protoc-gen-go-grpc.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-grpc-v1.1.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-go-grpc.mod -o=$(GOBIN)/protoc-gen-go-grpc-v1.1.0 "google.golang.org/grpc/cmd/protoc-gen-go-grpc"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.26.0
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.26.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.26.0 "google.golang.org/protobuf/cmd/protoc-gen-go"


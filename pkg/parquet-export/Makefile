.PHONY: all
all: proto

include .bingo/Variables.mk

.PHONY: proto
proto: $(BUF) $(PROTOC_GEN_GO_GRPC) $(PROTOC_GEN_GO)
	@cp $(PROTOC_GEN_GO) $(GOBIN)/protoc-gen-go
	@cp $(PROTOC_GEN_GO_GRPC) $(GOBIN)/protoc-gen-go-grpc
	@cd export1 && $(BUF) generate --template buf.gen.yaml .
	@cd export2 && $(BUF) generate --template buf.gen.yaml .

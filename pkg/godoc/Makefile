.PHONY: all
all: serve

include .bingo/Variables.mk

.PHONY: serve
serve: $(GODOC)
	$(GODOC) -http=:6060

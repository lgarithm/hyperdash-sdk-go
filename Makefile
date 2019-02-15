.PHONY: tools test

CMDS = \
	./cmd/...

tools:
	go version
	GOBIN=$(PWD)/bin \
	go install -v $(CMDS)

install:
	go install -v $(CMDS)

clean:
	go clean -v -cache ./...

test:
	go test -v ./...

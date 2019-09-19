DEFAULT_GOPATH=$(shell echo $$GOPATH|tr ':' '\n'|awk '!x[$$0]++'|sed '/^$$/d'|head -1)
ifeq ($(DEFAULT_GOPATH),)
DEFAULT_GOPATH := ~/go
endif
DEFAULT_GOBIN=$(DEFAULT_GOPATH)/bin
export PATH:=$(PATH):$(DEFAULT_GOBIN)

GOLANGCI_LINT=$(DEFAULT_GOBIN)/golangci-lint

default: clean build

deps:
	@GO111MODULE=on go mod vendor -v
	@GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go

build: bin/app bin/client

bin/%: cmd/%/main.go
	@GO111MODULE=on go build -o $@ $<

clean:
	@rm -rf bin/*

clean-all:
	@go clean --modcache
	
run: build
	@./bin/app

protoc-gen: api/*.pb.go

api/%.pb.go: api/%.proto 
	@protoc -I =api --go_out=plugins=grpc:api $<

$(GOLANGCI_LINT):
	@echo ">> Couldn't find $(GOLANGCI_LINT); installing"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | \
	sh -s -- -b $(DEFAULT_GOBIN) v1.15.0

show-linter:
	@echo $(GOLANGCI_LINT)

lint: $(GOLANGCI_LINT)
	@echo '>> Linting source code'
	@GO111MODULE=on $(GOLANGCI_LINT) \
	--enable=golint \
	--enable=gocritic \
	--enable=misspell \
	--enable=nakedret \
	--enable=unparam \
	--enable=lll \
	--enable=goconst \
	run ./...
VERSION = $(shell cat VERSION)
DVCS_HOST = github.com
ORG = geomyidia
PROJ = go-svc-conventions
FQ_PROJ = $(DVCS_HOST)/$(ORG)/$(PROJ)

LD_VERSION = -X $(FQ_PROJ)/pkg/version.version=$(VERSION)
LD_BUILDDATE = -X $(FQ_PROJ)/pkg/version.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_GITCOMMIT = -X $(FQ_PROJ)/pkg/version.gitCommit=$(shell git rev-parse --short HEAD)
LD_GITBRANCH = -X $(FQ_PROJ)/pkg/version.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
LD_GITSUMMARY = -X $(FQ_PROJ)/pkg/version.gitSummary=$(shell git describe --tags --dirty --always)
LDFLAGS = -w -s $(LD_VERSION) $(LD_BUILDDATE) $(LD_GITBRANCH) $(LD_GITSUMMARY) $(LD_GITCOMMIT)

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

build: protoc-regen bin/app bin/client

bin/%: cmd/%/main.go
	@GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o $@ $<

clean:
	@rm -rf bin/* 

clean-all:
	@go clean --modcache
	@rm -rf api/*.pb.go
	
run: build
	@./bin/app

protoc-gen: api/*.pb.go

protoc-regen: clean-protobuf
	@$(MAKE) protoc-gen

clean-protobuf:
	@rm -f api/*.pb.go

api/%.pb.go: api/%.proto 
	@protoc -I api --go_out=plugins=grpc:api $<
	@cp -r api/$(FQ_PROJ)/api/* api/
	@rm -rf api/$(DVCS_HOST)

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

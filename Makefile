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
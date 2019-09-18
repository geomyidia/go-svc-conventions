default: clean build

deps:
	GO111MODULE=on go mod vendor -v

build:
	GO111MODULE=on go build -o bin/app cmd/app/main.go

clean:
	rm -rf bin/*

run: build
	./bin/app

.PHONY: all cmdline web clean

GIT_COMMIT := $(shell git rev-parse --short HEAD)

all: web cmdline

cmdline:
	go build -ldflags "-X 'main.version=$(GIT_COMMIT)'" -o issue2md github.com/bigwhite/issue2md/cmd/issue2md

web:
	go build -o issue2mdweb

buildimage:
	docker build -t bigwhite/issue2mdweb .

push:
	docker push bigwhite/issue2mdweb:latest

clean:
	rm -fr issue2md issue2mdweb

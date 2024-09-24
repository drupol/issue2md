.PHONY: all cmdline web clean

all: web cmdline

cmdline:
	go build -o issue2md github.com/bigwhite/issue2md/cmd/issue2md

web:
	go build -o issue2mdweb 

clean:
	rm -fr issue2md issue2mdweb

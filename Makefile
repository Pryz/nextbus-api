BINARY=nextbus-api

test:
		go test -v

build:
		go get github.com/Pryz/nextbus
	  go get gopkg.in/redis.v5
		go build -o ${BINARY} *.go

.PHONY: clean
	clean:
		if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

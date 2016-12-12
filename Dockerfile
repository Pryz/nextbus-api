FROM golang
MAINTAINER Julien Fabre <ju.pryz@gmail.com>

RUN go get github.com/Pryz/nextbus-api

CMD [ "nextbus-api" ]

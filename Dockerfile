# syntax=docker/dockerfile:1

FROM golang:alpine 

WORKDIR /go/src/github.com/k0marov/go-socnet

COPY * ./

ENV GOPATH=/go

RUN go mod tidy
RUN go get -u ./...

RUN go build -o /go-socnet 

EXPOSE 8080 

CMD ["/go-socnet"] 

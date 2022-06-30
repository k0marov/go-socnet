# syntax=docker/dockerfile:1

FROM golang:alpine 

WORKDIR /app 

COPY * ./

ENV GOPATH=/app

RUN go mod tidy
RUN go get -u ./...

RUN go build -o /go-socnet 

EXPOSE 8080 

CMD ["/go-socnet"] 

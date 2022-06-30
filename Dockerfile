# syntax=docker/dockerfile:1

FROM golang:alpine 

WORKDIR /app 

COPY * ./ 

RUN go mod download 

RUN go build -o /go-socnet 

EXPOSE 8080 

CMD ["/go-socnet"] 

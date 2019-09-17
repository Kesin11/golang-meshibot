FROM golang:1.13-alpine

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -o /go/bin/golang-meshibot

CMD ["/go/bin/golang-meshibot"]
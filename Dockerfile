FROM golang:1.13-alpine
EXPOSE 8080

ENV GOPATH=/go
RUN mkdir -p $GOPATH/src/github.com/random-id
COPY . $GOPATH/src/github.com/random-id

WORKDIR $GOPATH/src/github.com/random-id
RUN go build -o randomid .

CMD ["/go/src/github.com/random-id/randomid"]

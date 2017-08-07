FROM golang:1.8.3

ADD . /go/src/github.com/ztplz/blog-server

RUN go install github.com/ztplz/blog-server

ENTRYPOINT /go/bin/blog-server

EXPOSE 8080
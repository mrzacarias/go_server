FROM golang

COPY . /go/src/github.com/mrzacarias/go_server

RUN go install github.com/mrzacarias/go_server

CMD ["/go/bin/go_server"]

EXPOSE 8080

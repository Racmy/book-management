FROM golang:latest

RUN go get -u github.com/oxequa/realize && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/go-sql-driver/mysql && \
    go get -u github.com/gorilla/sessions

WORKDIR /go/src/github.com/docker_go_nginx/

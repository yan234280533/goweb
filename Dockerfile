FROM golang:1.14 as builder

ADD . /go/src/go-web

WORKDIR /go/src/go-web

RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export GIT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags  -v -o main main.go && \
    mv main /go-web

FROM golang:1.14 as runtime

COPY --from=builder /go-web /go-web

CMD ["/go-web"]
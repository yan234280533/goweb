FROM debian:10-slim

COPY main /go-web

CMD ["/go-web"]
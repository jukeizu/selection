FROM golang:1.20 as build
WORKDIR /go/src/github.com/jukeizu/selection
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "nobody:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/jukeizu/selection/passwd /etc/passwd
COPY --from=build --chown=100:101 /go/src/github.com/jukeizu/selection/bin/selection .
USER nobody
ENTRYPOINT ["./selection"]

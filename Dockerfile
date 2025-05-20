FROM golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache make

WORKDIR /build
COPY main.go go.mod go.sum Makefile /build/
COPY pkg /build/pkg
COPY cmd /build/cmd

RUN make production-build CGO_ENABLED=0
RUN cp /build/bin/git-lost-and-found-$(go env GOOS)-$(go env GOARCH) /usr/local/bin/git-lost-and-found

FROM alpine:3.21

WORKDIR /target

RUN apk add --no-cache git

ENV PATH="$PATH:/usr/local/bin"
COPY --from=builder /usr/local/bin/git-lost-and-found /usr/local/bin/git-lost-and-found

ENTRYPOINT ["/usr/local/bin/git-lost-and-found"]

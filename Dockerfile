# Copyright (C) 2025 Bryce Thuilot
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the FSF, either version 3 of the License, or (at your option) any later version.
# See the LICENSE file in the root of this repository for full license text or
# visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

FROM golang:1.25-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS builder
ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS="-w -s"

WORKDIR /build
COPY main.go go.mod go.sum Makefile /build/
COPY pkg/ pkg/
COPY cmd/ cmd/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -ldflags="${LDFLAGS}" -trimpath -o git-lost-and-found main.go

FROM alpine/git:latest@sha256:bd54f921f6d803dfa3a4fe14b7defe36df1b71349a3e416547e333aa960f86e3

ENV PATH="$PATH:/usr/local/bin"
COPY --from=builder /build/git-lost-and-found /git-lost-and-found

ENTRYPOINT ["/git-lost-and-found"]

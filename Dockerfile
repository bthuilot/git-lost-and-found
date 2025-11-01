# Copyright (C) 2025 Bryce Thuilot
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the FSF, either version 3 of the License, or (at your option) any later version.
# See the LICENSE file in the root of this repository for full license text or
# visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

FROM golang:1.25-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS="-w -s"

WORKDIR /build
COPY main.go go.mod go.sum Makefile /build/
COPY pkg/ pkg/
COPY cmd/ cmd/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -ldflags="${LDFLAGS}" -trimpath -o git-lost-and-found main.go

FROM alpine/git:latest@sha256:94b40c2135951103e0c5c7db07ae4cf6e935644a717e05d17f0c540db47683af

ENV PATH="$PATH:/usr/local/bin"
COPY --from=builder /build/git-lost-and-found /git-lost-and-found

ENTRYPOINT ["/git-lost-and-found"]

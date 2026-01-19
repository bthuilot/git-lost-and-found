# Copyright (C) 2024-2026 Bryce Thuilot
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the FSF, either version 3 of the License, or (at your option) any later version.
# See the LICENSE file in the root of this repository for full license text or
# visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

FROM golang:1.25-alpine@sha256:d9b2e14101f27ec8d09674cd01186798d227bb0daec90e032aeb1cd22ac0f029 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS="-w -s"

WORKDIR /build
COPY main.go go.mod go.sum Makefile /build/
COPY pkg/ pkg/
COPY cmd/ cmd/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -ldflags="${LDFLAGS}" -trimpath -o git-lost-and-found main.go

FROM alpine/git:latest@sha256:d46d88ab234733c6b6a9771acd6d1384172fd0e2224e0232bdae32ec671aa099

ENV PATH="$PATH:/usr/local/bin"
COPY --from=builder /build/git-lost-and-found /git-lost-and-found

ENTRYPOINT ["/git-lost-and-found"]

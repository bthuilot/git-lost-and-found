FROM golang:1.22-bookworm

WORKDIR /build
RUN git clone --depth=1 https://github.com/gitleaks/gitleaks.git /build/gitleaks &&\
    cd  /build/gitleaks &&\
    make build &&\
    cp /build/gitleaks/gitleaks /bin/gitleaks

RUN curl -sSfL https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | sh -s -- -b /usr/local/bin

ENV PATH="$PATH:/bin"

WORKDIR /build
COPY main.go go.mod go.sum Makefile /build/
COPY pkg /build/pkg
COPY cmd /build/cmd

RUN make build && cp /build/bin/git-lost-and-found /bin/git-lost-and-found

ENTRYPOINT ["/bin/git-lost-and-found"]

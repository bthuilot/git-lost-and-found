FROM golang:1.22-bookworm

WORKDIR /build
RUN git clone --depth=1 https://github.com/gitleaks/gitleaks.git /build/gitleaks &&\
    cd  /build/gitleaks &&\
    make build &&\
    cp /build/gitleaks/gitleaks /bin/gitleaks

ENV PATH="$PATH:/bin"

WORKDIR /build/git-scanner
COPY main.go go.mod go.sum Makefile /build/git-scanner/
COPY pkg /build/git-scanner/pkg
COPY cmd /build/git-scanner/cmd

RUN make build && cp /build/git-scanner/bin/git-scanner /bin/git-scanner

ENTRYPOINT ["/bin/git-scanner"]

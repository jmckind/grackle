# Golang image for building the binary
FROM golang:alpine AS builder

# Copy artifacts
WORKDIR /go/src/github.com/jmckind/grackle
COPY ./ ./

# Compile go binaries
RUN set -x && \ 
    apk add --no-cache git ca-certificates && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o grackle-ingest cmd/ingest/main.go && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o grackle-web cmd/web/main.go && \
    ls -lh ./grackle-*

# Universal base image for release image
FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

# Setup environment and working directory
ENV USER_UID=1001 \
    USER_NAME=grackle \
    USER_HOME=/opt/grackle
WORKDIR /

# Copy builder artifacts
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/jmckind/grackle/grackle-* /

# Copy static artifacts
COPY bin/entrypoint.sh /entrypoint
COPY templates/* /templates/

# Setup application user
RUN mkdir -p ${USER_HOME} && \
    chown ${USER_UID}:0 ${USER_HOME} && \
    chmod ug+rwx ${USER_HOME} && \
    chmod g+rw /etc/passwd

# Set entrypoint 
ENTRYPOINT ["/entrypoint"]

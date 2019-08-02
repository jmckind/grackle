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

# Minimal image for final binary
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/jmckind/grackle/grackle-* /
COPY templates/* /templates/

WORKDIR /
CMD [ "/grackle-web" ]

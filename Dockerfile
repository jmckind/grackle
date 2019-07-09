# Golang image for building the binary
FROM golang:alpine AS builder

WORKDIR /go/src/github.com/jmckind/grackle
COPY ./ ./

RUN set -x && \ 
    apk add git --no-cache && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o grackle-ingest cmd/grackle-ingest/main.go && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o grackle-web cmd/grackle-web/main.go && \
    ls -lh ./grackle-*

# Minimal image for final binary
FROM scratch

COPY --from=builder /go/src/github.com/jmckind/grackle/grackle-* /
COPY templates/* /templates/

WORKDIR /
CMD [ "/grackle-web" ]

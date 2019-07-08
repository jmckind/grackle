build: build-ingest build-web

build-ingest:
	mkdir -p build
	go build -a -o build/grackle-ingest cmd/grackle-ingest/main.go

build-web:
	mkdir -p build
	go build -a -o build/grackle-web cmd/grackle-web/main.go

clean:
	rm -fr build

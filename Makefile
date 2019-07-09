build: build-ingest build-web

build-ingest:
	mkdir -p build
	go build -a -o build/grackle-ingest cmd/grackle-ingest/main.go

build-web:
	mkdir -p build
	go build -a -o build/grackle-web cmd/grackle-web/main.go

clean:
	rm -fr build

docker:
	docker build -t jmckind/grackle:latest .

run-rethinkdb:
	docker run -d -p 28015:28015 -p 29015:29015 -p 8080:8080 rethinkdb:latest

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

run-docker-grackle-web:
	docker run -it jmckind/grackle:latest /grackle-web --rethinkdb-host 172.17.0.1

run-local-grackle-ingest:
	go run cmd/grackle-ingest/main.go

run-local-grackle-web:
	go run cmd/grackle-web/main.go

run-local-rethinkdb:
	docker run -d -p 28015:28015 -p 29015:29015 -p 8080:8080 rethinkdb:latest

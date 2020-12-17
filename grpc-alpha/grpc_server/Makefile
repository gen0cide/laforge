GRPCPORT = ":50051"
GINPORT = ":5000"

build:
	go build -ldflags="-X 'main.port=$(GRPCPORT)' -X 'main.webPort=$(GINPORT)' " -o laforge_server ./...

run:
	go run -ldflags="-X 'main.port=$(GRPCPORT)' -X 'main.webPort=$(GINPORT)' " ./...

clean:
	rm laforge_server *.db

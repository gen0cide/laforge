GRPCPORT = ":50051"
GINPORT = ":5000"

copy_certs:
	cp ../certs/service.key ./certs/
	cp ../certs/service.pem ./certs/

gen_static:
	fileb0x assets.toml

build: copy_certs gen_static
	go build -ldflags="-X 'main.port=$(GRPCPORT)' -X 'main.webPort=$(GINPORT)' " -o laforge_server github.com/gen0cide/laforge/grpc/server

run: copy_certs gen_static
	go run -ldflags="-X 'main.port=$(GRPCPORT)' -X 'main.webPort=$(GINPORT)' " github.com/gen0cide/laforge/grpc/server

clean:
	rm laforge_server *.db ./certs/* 
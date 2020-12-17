GPRCADDRESS = "localhost:50051"
AgentID 	= "45846255432156"

copy_certs:
	cp ../grpc_certs/service.pem ./certs/

gen_static:
	fileb0x assets.toml

build: copy_certs gen_static
	go build -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" -o laforge_agent github.com/gen0cide/laforge/grpc-alpha/grpc_client

run: copy_certs gen_static
	go run -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" github.com/gen0cide/laforge/grpc-alpha/grpc_client

clean:
	rm laforge_agent ./certs/* 
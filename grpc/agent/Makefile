GPRCADDRESS = "18.232.181.162:50051"
AgentID 	= "3"

copy_certs:
	cp ../certs/service.pem ./certs/

gen_static:
	fileb0x assets.toml

build: copy_certs gen_static
	go build -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" -o laforge_agent github.com/gen0cide/laforge/grpc/agent

run:
	go run -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" github.com/gen0cide/laforge/grpc/agent

clean:
	rm laforge_agent ./certs/* 
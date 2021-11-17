GPRCADDRESS = "localhost:50051"
AgentID 	= "3"

copy_certs:
	cp ../certs/service.pem ./certs/

gen_static:
	fileb0x assets.toml

build_linux: copy_certs gen_static
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" -o laforge_agent github.com/gen0cide/laforge/grpc/agent

build_windows: 
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zcc" go build -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" -o laforge_agent.exe github.com/gen0cide/laforge/grpc/agent

run:
	go run -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" github.com/gen0cide/laforge/grpc/agent

clean:
	rm laforge_agent ./certs/* 

GPRCADDRESS = "localhost:50051"
AgentID 	= "45846255432156"

build:
	go build -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" -o laforge_agent ./...

run:
	go run -ldflags=" -X 'main.clientID=$(AgentID)' -X 'main.address=$(GPRCADDRESS)'" ./...

clean:
	rm laforge_agent
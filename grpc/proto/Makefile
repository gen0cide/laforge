build: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative agent_proto.proto

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	
clean:
	rm *.go
## TESTING REDIS WITH GO ##
FROM golang:1.14

WORKDIR /go/src/app
COPY ./grpc_server .
COPY ./laforge_proto_agent /go/src/github.com/frybin/laforge/grpc-alpha/laforge_proto_agent

RUN go get -d -v ./...
RUN go build ./server.go

CMD ["./server"]
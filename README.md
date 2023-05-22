# go-grpc
go-grpc

protobuf https://github.com/protocolbuffers/protobuf
go-grpc_out https://github.com/golang/protobuf/tree/master/protoc-gen-go

```shell
protoc --go_out=. --go-grpc_out=. ./simple/*.proto
protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. ./simple/*.proto
```
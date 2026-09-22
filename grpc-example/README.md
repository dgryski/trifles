This is a sample grpc backend for the proxy example.

To test, first install `grpcurl`:
```
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

To build and run the server:
```
$ go run main.go
2026/09/22 12:46:07 Serving gRPC on 0.0.0.0:9999
```

To send a request:
```
$ grpcurl -proto hello.proto  -plaintext localhost:9999 helloworld.Greeter/SayHello
{
  "message": " world"
}
$ grpcurl -d '{"name": "Fastly"}' -proto hello.proto  -plaintext localhost:9999 helloworld.Greeter/SayHello
{
  "message": "Fastly world"
}
```

If running through the `grpc-proxy` example, the port changes to `:7676`
```
$ grpcurl -d '{"name": "Fastly"}' -proto hello.proto  -plaintext localhost:7676 helloworld.Greeter/SayHello
{
  "message": "Fastly world"
}
```

To regenerate the service definitions:
```
$ brew install protobuf buf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$ buf generate
```

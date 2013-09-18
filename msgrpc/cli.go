package main

import (
	"fmt"
	codec "github.com/ugorji/go/codec"
	"log"
	"net"
	"net/rpc"
)

// dup'd in srv.go
type EchoArgument struct {
	A string
	B string
	C string
}

func main() {

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var h codec.MsgpackHandle
	rpcCodec := codec.MsgpackSpecRpc.ClientCodec(conn, &h)
	client := rpc.NewClientWithCodec(rpcCodec)

	args := []string{"foo1", "foo2", "foo3"}
	var reply string
	err = client.Call("Echo123", args, &reply)
	if err != nil {
		log.Fatal("Echo123 error:", err)
	}
	fmt.Println("reply=", reply)

	arg := &EchoArgument{A: "a string", B: "b string", C: "c string"}
	err = client.Call("EchoStruct", arg, &reply)
	if err != nil {
		log.Fatal("EchoStruct error:", err)
	}
	fmt.Println("reply=", reply)

}

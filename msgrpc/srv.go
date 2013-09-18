package main

import (
	"fmt"
	codec "github.com/ugorji/go/codec"
	"log"
	"net"
	"net/rpc"
)

type Echo struct{}

func (e *Echo) Echo123(args []string, response *string) error {
	log.Printf("args=%#v\n", args)
	*response = fmt.Sprintf("%#v", args)
	return nil
}

// dup'd in cli.go
type EchoArgument struct {
	A string
	B string
	C string
}

func (e *Echo) EchoStruct(arg EchoArgument, response *string) error {
	log.Printf("args=%#v\n", arg)
	*response = fmt.Sprintf("%#v", arg)
	return nil
}

func main() {

	echo := &Echo{}
	rpc.Register(echo)

	ln, e := net.Listen("tcp", ":9000")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	log.Println("rpc server starting")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		var h codec.MsgpackHandle
		c := codec.MsgpackSpecRpc.ServerCodec(conn, &h)

		go rpc.ServeCodec(c)
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()
	rawurl := flag.Arg(0)

	// ws://echo.websocket.org:80
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal("error parsing: ", err)
	}

	conn, err := websocket.Dial(rawurl, "tcp", "http://"+u.Host)

	if err != nil {
		log.Fatal("error dialing: ", err)
	}

	go func(conn *websocket.Conn) {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			s := scanner.Text()
			fmt.Println("send=", s)
			websocket.Message.Send(conn, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("error scanning stdin: %s", err)
		}
	}(conn)

	for {
		var msg string
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Fatal("error recv: ", err)
		}
		fmt.Println("recv=", msg)
	}
}

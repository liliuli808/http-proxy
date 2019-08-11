package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(conn net.Conn) {
	if conn == nil {
		log.Println("no client")
		return
	}
	defer conn.Close()

	var b [1024]byte

	n, err := conn.Read(b[:])

	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)

	hostPort, err := url.Parse(host)

	if err != nil {
		log.Println(err)
		return
	}

	if hostPort.Opaque == "443" {
		address = hostPort.Scheme + ":443"
	} else {
		if strings.Index(hostPort.Host, ":") == -1 {
			address = hostPort.Host + "80"
		} else {
			address = hostPort.Host
		}
	}

	server, err := net.Dial("tcp", address)

	if err != nil {
		log.Println(err)
		return
	}

	if method == "CONNECT" {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}

	go io.Copy(server, conn)
	io.Copy(conn, server)
}

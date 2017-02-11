package main

import (
	"io"
	"log"
	"net"
	"time"
)

var transportBufferSize = 32 * 1024

var payload = "zazkia shamelessly says super stuff"

func main() {
	go startServer()
	startClientLoop()
}

func startClientLoop() {
	for {
		log.Println("client dialing :54321 ...")
		serviceConn, err := net.Dial("tcp", "localhost:54321")
		if err != nil {
			log.Fatalf("failed to connect to server because [%v]", err)
		}
		log.Println("client connected to :54321")
		for {
			time.Sleep(1 * time.Second)
			log.Println("sending", payload)
			_, err := io.WriteString(serviceConn, payload)
			if err != nil {
				log.Printf("failed to send message over connections because [%v]", err)
				break
			}
		}
		serviceConn.Close()
	}
}

func startServer() {
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("failed to start listener:%v", err)
	}
	log.Println("server accepting connections on :12345")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept new connections because [%v]", err)
			break
		}
		go handleConnection(conn)
	}
}

func handleConnection(clientConn net.Conn) {
	for {
		buffer := make([]byte, transportBufferSize)
		read, err := clientConn.Read(buffer)
		if err != nil {
			log.Printf("failed to read from connection because [%v]", err)
			break
		}
		log.Printf("%d bytes read : %s\n", read, string(buffer[:read]))
	}
}

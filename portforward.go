package main

import (
	"io"
	"os"
	"net"
	"log"
	"math/rand"
	"time"
)

func main() {
	if len(os.Args) <= 3 {
		log.Fatal("usage: portforward local remote1 remote2 remote3 ...")
	}
	localAddrString := os.Args[1]
	remoteAddrStrings := os.Args[2:]
	localAddr, err := net.ResolveTCPAddr("tcp", localAddrString)
	if localAddr == nil {
		log.Fatalf("net.ResolveTCPAddr failed: %s", err)
	}
	local, err := net.ListenTCP("tcp", localAddr)
	if local == nil {
		log.Fatalf("portforward: %s", err)
	}
	log.Printf("portforward listen on %s", localAddr)

	rand.Seed(time.Now().UnixNano())
	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Fatalf("accept failed: %s", err)
		}
		go forward(conn, remoteAddrStrings)
	}
}

func forward(local net.Conn, remoteAddrs []string) {
	remoteAddr := remoteAddrs[rand.Intn(len(remoteAddrs))]
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		local.Close()
		return
	}
	go func() {
		defer local.Close()
		//remote.SetReadTimeout(120*1E9)
		io.Copy(local, remote)
	}()
	go func() {
		defer remote.Close()
		//local.SetReadTimeout(120*1E9)
		io.Copy(remote, local)
	}()
	log.Printf("forward %s to %s", local.RemoteAddr(), remoteAddr)
}

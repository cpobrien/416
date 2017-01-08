package main

import (
	"fmt"
	"os"
	"net"
	"time"
	"bytes"
	"encoding/gob"
	"bufio"
	"log"
)

type Network struct {
	conn   net.Conn
	reader *bufio.Reader
}

var BufferSize int = 1024

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}

func (a *Network) send(i int) string {
	var buffer []byte
	for {
		payload, _ := Marshall(uint32(i))
		_, err := a.conn.Write(payload)
		if err == nil {
			break
		}
		fmt.Println("not snet.")
	}
	for {
		_, err := a.conn.Read(buffer)
		if err == nil {
			break
		}
	}
	return string(buffer)
}

func (a *Network) Run() {
	for {
		fmt.Println("Input a number.")
		var i int
		var err error
		for {
			i , err = fmt.Scanf("%d", &i)
			fmt.Println("read")
			if err != nil {
				fmt.Print("Invalid number.\nNumber > ")
			} else {
				break
			}
		}
		res := a.send(i)
		switch res {
		case "high": fallthrough
		case "low": fmt.Println("Your number is too %s.", res)
		default:
			fmt.Println("%s", res)
			a.conn.Close()
			return
		}
	}
}

func NewAssignment(local, remote string) *Network {
	a := new(Network)
	raddr, _ := net.ResolveUDPAddr("udp", remote)
	laddr, _ := net.ResolveUDPAddr("udp", local)
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	conn.SetReadBuffer(BufferSize)
	a.conn = conn
	return a
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: client.go [local UDP ip:port] [server UDP ip:port]")
		return
	}
	local_ip_port := args[0]
	remote_ip_port := args[1]
	NewAssignment(local_ip_port, remote_ip_port).Run()
}
package main

import (
	"fmt"
	"os"
	"net"
	"bytes"
	"encoding/gob"
	"bufio"
	"log"
	"time"
)

type Network struct {
	laddr  *net.UDPAddr
	raddr  *net.UDPAddr
	reader *bufio.Reader
}

func (n *Network)StartUDP() *net.UDPConn {
	conn, err := net.DialUDP("udp", n.laddr, n.raddr)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Second))
	return conn
}

var BufferSize int = 1024

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}

func (a *Network) send(i int) string {
	buf := make([]byte, BufferSize)
	payload, _ := Marshall(uint32(i))
	for {
		conn := a.StartUDP()
		_, err := conn.Write(payload)
		if err != nil {
			log.Fatal(err)
		}
		l, err := conn.Read(buf)
		conn.Close()
		// If at first read not succeed, try try again
		if l != 0 || err == nil {
			break
		}
		return string(buf[:l])
	}
}

func (a *Network) Run() {
	for {
		fmt.Println("Input a number.")
		var i int
		var err error
		for {
			_, err = fmt.Scanf("%d", &i)
			if err != nil {
				fmt.Print("Invalid number.\nNumber > ")
			} else {
				break
			}
		}
		res := a.send(i)
		switch res {
		case "high": fallthrough
		case "low": fmt.Printf("Your number is too %s.\n", res)
		default:
			fmt.Println(res)
			return
		}
	}
}

func NewNetwork(local, remote string) *Network {
	a := new(Network)
	raddr, _ := net.ResolveUDPAddr("udp", remote)
	laddr, _ := net.ResolveUDPAddr("udp", local)
	a.raddr = raddr
	a.laddr = laddr
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
	NewNetwork(local_ip_port, remote_ip_port).Run()
}
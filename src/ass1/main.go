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
	local  *net.UDPAddr
	remote *net.UDPAddr
	reader *bufio.Reader
}

func (n *Network)StartUDP() *net.UDPConn {
	conn, err := net.DialUDP("udp", n.local, n.remote)
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

func (n *Network) send(i int) string {
	var len int
	buf := make([]byte, BufferSize)
	payload, _ := Marshall(uint32(i))
	for {
		conn := n.StartUDP()
		_, err := conn.Write(payload)
		if err != nil {
			log.Fatal(err)
		}
		len, err = conn.Read(buf)
		conn.Close()
		// If read not succeed, try try again
		if len != 0 && err == nil {
			break
		}
	}
	return string(buf[:len])
}

func (n *Network) Run() {
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
		res := n.send(i)
		switch res {
		case "high", "low": fmt.Printf("Your number is too %s.\n", res)
		default:
			fmt.Println(res)
			return
		}
	}
}

func NewNetwork(local, remote string) *Network {
	n := new(Network)
	n.remote, _ = net.ResolveUDPAddr("udp", remote)
	n.local, _ = net.ResolveUDPAddr("udp", local)
	return n
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
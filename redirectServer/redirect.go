package redirectServer

import (
	"fmt"
	"net"
)

type Redirect struct {
}

func (r *Redirect) Run(addrA, addrB string) {
	a, err := net.ResolveTCPAddr("tcp", addrA)
	if err != nil {
		fmt.Println(err)
	}
	b, err := net.ResolveTCPAddr("tcp", addrB)
	if err != nil {
		fmt.Println(err)
	}
	r.connection(a, b)

}

func (r *Redirect) connection(serverA, serverB *net.TCPAddr) {
	connA, err := net.DialTCP("tcp", nil, serverA)
	if err != nil {
		fmt.Println(err)
	}
	connB, err := net.DialTCP("tcp", nil, serverB)
	if err != nil {
		fmt.Println(err)
	}
	go func(a, b *net.TCPConn) {
		for {
			buffer := make([]byte, 1024)
			n, err := a.Read(buffer)
			if err != nil {
				fmt.Println(err)
			} else {
				b.Write(buffer[:n])
				fmt.Printf("a ====> b : %X", buffer[:n])
			}
		}
	}(connA, connB)
	go func(a, b *net.TCPConn) {
		for {
			buffer := make([]byte, 1024)
			n, err := b.Read(buffer)
			if err != nil {
				fmt.Println(err)
			} else {
				a.Write(buffer[:n])
				fmt.Printf("b ====> a : %X", buffer[:n])
			}
		}
	}(connA, connB)
}

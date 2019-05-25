package redirectServer

import (
	"fmt"
	"io"
	"net"
	"time"
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
	go func(from, to *net.TCPConn) {
		for {
			n, err := io.Copy(from, to)
			if err != nil {
				fmt.Println(err)
				<-time.After(time.Second)
			} else {
				fmt.Printf("a ====> b : 转发 %d 个字节 \n", n)
			}
		}
	}(connA, connB)
	go func(from, to *net.TCPConn) {
		for {
			n, err := io.Copy(from, to)
			if err != nil {
				fmt.Println(err)
				<-time.After(time.Second)
			} else {
				fmt.Printf("b ====> a : 转发 %d 个字节 \n", n)
			}
		}
	}(connB, connA)
}

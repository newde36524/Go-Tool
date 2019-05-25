package redirectserver

import (
	"fmt"
	"io"
	"net"
)

//Redirect .
type Redirect struct {
}

//Run .
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

//connection .
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
		n, err := io.Copy(from, to)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("a ====> b : 转发 %d 个字节 \n", n)
		}
	}(connA, connB)
	go func(from, to *net.TCPConn) {
		n, err := io.Copy(from, to)
		if err != nil {
			fmt.Println("错误信息：", err)
		} else {
			fmt.Printf("b ====> a : 转发 %d 个字节 \n", n)
		}
	}(connB, connA)
}

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/christianwoehrle/tcpidaler/person"
	"github.com/golang/protobuf/proto"
)

func main() {

	fmt.Printf("start: %d\n", 2018)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1201")
	//listener, err := net.ListenTCP("tcp", tcpAddr)
	wg := sync.WaitGroup{}
	wg.Add(1)
	quit := make(chan interface{})
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Printf("Created Listener with errcode %v\n", err)
	go func() {
		for {
			select {
			case <-quit:

				wg.Done()
				return

			default:

			}
			pn := person.Person_Name{Family: "wöhrle", Personal: "pers"}
			pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
			pes := []*person.Person_Email{&pe}
			p := person.Person{Name: &pn, Email: pes}

			out, err := proto.Marshal(&p)
			fmt.Println("protobuf", out)
			listener.SetDeadline(time.Now().Add(2 * time.Second))
			lconn, err := listener.AcceptTCP()
			fmt.Println("Accept TCP, err: ", err)
			i, err := lconn.Write([]byte("Dude"))
			fmt.Printf("Written %d bytes with errcode %v\n", i, err)
			lconn.Close()
		}
	}()

	port, err := net.LookupPort("tcp", "smtp")
	fmt.Println("Service", port)
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "localhost:1201")
	fmt.Println(err)
	conn, err := net.DialTCP("tcp4", nil, tcpaddr)
	fmt.Println(err)
	i, err := conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	fmt.Println(i, err)
	read, _ := ioutil.ReadAll(conn)
	fmt.Println(string(read))

	wg.Wait()
}

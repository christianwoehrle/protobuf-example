package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"os"

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
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	handleErr("Error creating Server Listener", err)
	go server(*listener, quit)
	go client()
	go client()
	go client()
	time.Sleep(3 * time.Second)
	wg.Wait()
}

func client() {
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "localhost:1201")
	handleErr("Couldn´´ resolve Address: ", err)
	conn, err := net.DialTCP("tcp4", nil, tcpaddr)
	defer conn.Close()
	fmt.Println(err)
	pn := person.Person_Name{Family: "wöhrle", Personal: "pers"}
	pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
	pes := []*person.Person_Email{&pe}
	p := person.Person{Name: &pn, Email: pes}

	out, err := proto.Marshal(&p)
	fmt.Println("protobuf", out)
	i, err := conn.Write([]byte(out))
	fmt.Println("protobuf rausgeschrieben", i, err)

}

func server(listener net.TCPListener, quit <-chan interface{}) {

	for {

		select {
		case <-quit:
			return

		default:

		}
		lConn, err := listener.AcceptTCP()
		handleErr("Got connection", err)
		var read [512]byte

		n, err := lConn.Read(read[0:])

		//read, err := ioutil.ReadAll(lconn)
		fmt.Println("ReadProtobuf from Server, #bytes: ", n, read)
		pb := person.Person{}
		proto.Unmarshal(read[0:], &pb)
		fmt.Printf("PB: %v\n", pb)
		//lConn.Close()
	}
}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}

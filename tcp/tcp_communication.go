package main

// Simple program that uses protobuf strctures but not the protobuf messages
// tcp listener is set up, accepting connection request and handling the data messages (just echoing)
//

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/christianwoehrle/protobuf-example/person"
	"github.com/gogo/protobuf/proto"
)

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8889")
	handleErr("Cannot Resove TCPAddr", err)
	serverReadyWG := sync.WaitGroup{}
	serverReadyWG.Add(1)
	quit := make(chan interface{})

	go serverTcp(tcpAddr, quit, &serverReadyWG)

	serverReadyWG.Wait()

	clientDoneWg := sync.WaitGroup{}
	for i := 1; i < 100; i++ {
		fmt.Println(i)
		clientDoneWg.Add(1)
		go clientTcp(&clientDoneWg)

	}
	clientDoneWg.Wait()
	time.Sleep(1 * time.Second)
	close(quit)
}

func clientTcp(clientDoneWg *sync.WaitGroup) {
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "localhost:1201")
	handleErr("Couldn´t resolve Address: ", err)
	conn, err := net.DialTCP("tcp4", nil, tcpaddr)
	handleErr("Cannot Dial Server", err)
	defer conn.Close()
	pn := person.Person_Name{Family: "woehrle", Personal: "pers"}
	pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
	pes := []*person.Person_Email{&pe}
	p := person.Person{Name: &pn, Email: pes}
	out, err := proto.Marshal(&p)
	handleErr("Client Cannot Marshal Person", err)
	conn.Write(out)
	handleErr("Client Cannot Write to Server", err)

	var read [512]byte
	n, err := conn.Read(read[0:])
	handleErr("Client can´t read Response", err)
	//fmt.Println(read)
	pb := person.Person{}
	err = proto.Unmarshal(read[0:n], &pb)
	handleErr("Client can´t Unmarshal", err)
	conn.Close()
	clientDoneWg.Done()

}

func serverTcp(tcpAddr *net.TCPAddr, quit <-chan interface{}, serverReadyWg *sync.WaitGroup) {
	fmt.Println("called serverTcp")
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	handleErr("Error creating Server Listener", err)
	serverReadyWg.Done()
	for {

		select {
		case <-quit:
			fmt.Println("quit")
			return

		default:

		}
		lConn, err := listener.AcceptTCP()
		handleErr("Got connection", err)

		var read [512]byte
		n, err := lConn.Read(read[0:])

		//fmt.Println("Server Read: ", n, read)
		handleErr("Cannot Read Answer", err)
		pb := person.Person{}
		proto.Unmarshal(read[0:n], &pb)

		out, err := proto.Marshal(&pb)
		handleErr("Cannot Marshal Person", err)
		lConn.Write(out)
		handleErr("Cannot Write to Client", err)
		lConn.Close()
	}
}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}

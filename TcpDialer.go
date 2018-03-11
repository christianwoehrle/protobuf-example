package main

import (
	"fmt"
	"net"
	"sync"

	"os"

	"time"

	"github.com/christianwoehrle/tcpidaler/person"
	"github.com/golang/protobuf/proto"
)

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1201")
	handleErr("Cannot Resove TCPAddr", err)
	serverReadyWG := sync.WaitGroup{}
	serverReadyWG.Add(1)
	quit := make(chan interface{})

	go server(tcpAddr, quit, &serverReadyWG)

	serverReadyWG.Wait()

	fmt.Printf("server ready passed")
	clientDoneWg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		clientDoneWg.Add(1)
		go client(&clientDoneWg)

	}
	clientDoneWg.Wait()
	close(quit)
	time.Sleep(1 * time.Second)
	fmt.Println("und raus")
}

func client(clientDoneWg *sync.WaitGroup) {
	tcpaddr, err := net.ResolveTCPAddr("tcp4", "localhost:1201")
	handleErr("Couldn´t resolve Address: ", err)
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
	clientDoneWg.Done()

}

func server(tcpAddr *net.TCPAddr, quit <-chan interface{}, serverReadyWg *sync.WaitGroup) {
	fmt.Printf("start server")
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	handleErr("Error creating Server Listener", err)
	fmt.Printf("server ready")
	serverReadyWg.Done()
	for {

		select {
		case <-quit:
			fmt.Println("quit")
			return

		default:
			fmt.Printf("default")
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

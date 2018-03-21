package main

// Simple program that uses protobuf strctures but not the protobuf messages
// tcp listener is set up, accepting connection request and handling the data messages (just echoing)
//

import (
	"fmt"
	"net"

	"os"

	"time"

	"github.com/christianwoehrle/protobuf-example/person"
	"github.com/prometheus/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type PersonService struct{}

func (p PersonService) Echo(context context.Context, person *person.Person) (*person.Person, error) {
	return person, nil
}

func main() {
	go server()
	time.Sleep(1 * time.Second)
	clientEcho()

}

func clientEcho() {
	clientConnection, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	handleErr("Could not Dial grpc", err)
	client := person.NewPersonServiceClient(clientConnection)

	pn := person.Person_Name{Family: "woehrle", Personal: "pers"}
	pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
	pes := []*person.Person_Email{&pe}
	p := person.Person{Name: &pn, Email: pes}
	person2, err := client.Echo(context.Background(), &p, grpc.FailFast(true))
	fmt.Println(person2, err)
}

func server() {
	srv := grpc.NewServer()
	var pss PersonService
	person.RegisterPersonServiceServer(srv, pss)
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	handleErr("Could not Resolve Addr", err)
	listener, _ := net.ListenTCP("tcp", addr)
	handleErr("Could not create Listener", err)
	log.Fatal(srv.Serve(listener))
}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}

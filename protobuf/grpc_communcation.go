package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/christianwoehrle/protobuf-example/person"
	"google.golang.org/grpc"
)

// Simple program that uses protobuf strctures but not the protobuf messages
// tcp listener is set up, accepting connection request and handling the data messages (just echoing)
//

/**
  PersonService offers the Methods Echo and GetPeronStream via grpc
*/
type PersonService struct{}

/**
  Echo takes a Person as an argument and copies it back to the sender
*/
func (p PersonService) Echo(ctx context.Context, person *person.Person) (*person.Person, error) {
	return person, nil
}

/**
  GetPersonStream sends some Person Objects back to the Caller via a Stream
*/
func (p PersonService) GetPersonStream(e *person.Empty, stream person.PersonService_GetPersonStreamServer) error {

	for i := 0; i < 10; i++ {
		pn := person.Person_Name{Family: "woehrle" + strconv.Itoa(i), Personal: "pers"}
		pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
		pes := []*person.Person_Email{&pe}
		p := person.Person{Name: &pn, Email: pes}
		stream.Send(&p)
	}
	return nil
}

func main() {
	go server()
	time.Sleep(1 * time.Second)
	clientEcho()
	getPersonStream()

}

func getPersonStream() error {
	clientConnection, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	handleErr("Could not Dial grpc", err)
	client := person.NewPersonServiceClient(clientConnection)

	stream, err := client.GetPersonStream(context.Background(), &person.Empty{})
	handleErr("Couldn´´ call GetPersonStream", err)
	for {
		p, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		fmt.Println(p)

	}
	return nil

}

func clientEcho() (*person.Person, error) {

	pn := person.Person_Name{Family: "woehrle", Personal: "pers"}
	pe := person.Person_Email{Kind: "job", Address: "cw@gm.com"}
	pes := []*person.Person_Email{&pe}
	p := person.Person{Name: &pn, Email: pes}
	return clientEcho2(&p)
}

func clientEcho2(p *person.Person) (*person.Person, error) {
	clientConnection, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	handleErr("Could not Dial grpc", err)
	client := person.NewPersonServiceClient(clientConnection)

	person2, err := client.Echo(context.Background(), p, grpc.FailFast(true))
	return person2, err
}

func server() {
	srv := grpc.NewServer()
	var pss PersonService
	person.RegisterPersonServiceServer(srv, pss)
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	handleErr("Could not Resolve Addr", err)
	listener, _ := net.ListenTCP("tcp", addr)
	handleErr("Could not create Listener", err)
	log.Fatalln(srv.Serve(listener))
}

func handleErr(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", text, err)
		os.Exit(1)
	}
}

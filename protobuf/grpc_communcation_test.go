package main

import (
	"testing"

	"fmt"

	"time"

	"github.com/christianwoehrle/protobuf-example/person"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
)

//var update = flag.Bool("update", false, "update golden files")

func TestGrpcCall(t *testing.T) {

	go server()

	p, err := clientEcho()
	//fmt.Println(p)

	assert.Nil(t, err, fmt.Sprint("error when issuing grpc call", p, err))
}

func TestMain(m *testing.M) {
	fmt.Println("First--------------------------- TestMain setup")
	os.Exit(m.Run())
}
func TestGrpcCallTable(t *testing.T) {
	t.Helper()

	pp := []person.Person{{Name: &person.Person_Name{Family: "woehrle", Personal: "pers"}, Email: []*person.Person_Email{{Kind: "job", Address: "cw@gm.com"}}},
		{Name: &person.Person_Name{Family: "guenther", Personal: "pers"}, Email: []*person.Person_Email{{Kind: "home", Address: "cw@gmailm.com"}}}}

	for i, p := range pp {
		t.Run(strconv.Itoa(i)+"_"+p.Name.Family, func(t *testing.T) {
			result, err := clientEcho2(&p)
			assert.Nil(t, err, fmt.Sprint("error when issuing grpc call", result, err))
			assert.EqualValues(t, &p, result, "Echoed Person not the same")

		})
	}
}

func BenchmarkGrpcCall(b *testing.B) {

	go server()
	time.Sleep(1 * time.Second)
	b.Run("BenchmarkGrpcCallB.N", func(b *testing.B) {

		b.StartTimer()
		fmt.Println(b.N)
		for i := 1; i < b.N; i++ {
			p, err := clientEcho()
			if err != nil {
				b.Error("Call failed", err, p)
			}

		}
		b.StopTimer()
	})
	/*b.Run("BenchmarkGrpcCall100", func(b *testing.B) {
		b.StartTimer()
		for i := 1; i < 100; i++ {
			p, err := clientEcho()
			if err != nil {
				b.Error("Call failed", err, p)
			}

		}
		b.StopTimer()
	})*/

}

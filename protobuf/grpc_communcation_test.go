package main

import (
	"testing"

	"fmt"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestGrpcCall(t *testing.T) {

	go server()

	p, err := clientEcho()
	//fmt.Println(p)

	assert.Nil(t, err, fmt.Sprint("error when issuing grpc call", p, err))
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

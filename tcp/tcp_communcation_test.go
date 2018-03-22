package main

import (
	"net"
	"sync"
	"testing"
)

func BenchmarkTcpCall(b *testing.B) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1201")
	handleErr("Cannot Resove TCPAddr", err)
	serverReadyWG := sync.WaitGroup{}
	serverReadyWG.Add(1)
	quit := make(chan interface{})

	go serverTcp(tcpAddr, quit, &serverReadyWG)

	serverReadyWG.Wait()

	clientDoneWg := sync.WaitGroup{}
	b.Run("BenchmarkGrpcCallB.N", func(b *testing.B) {
		b.StartTimer()
		for i := 1; i < b.N; i++ {
			clientDoneWg.Add(1)
			clientTcp(&clientDoneWg)

		}
		b.StopTimer()
	})
	close(quit)

}

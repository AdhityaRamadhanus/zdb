package main

import (
	"github.com/AdhityaRamadhanus/zdb/tcp"
)

func main() {
	srv := tcp.NewServer("tcp", "localhost:9000")
	srv.Run()
}

package main

import (
	"flag"
	"fmt"

	"github.com/aeye/watcher-core/internal/server"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	fmt.Println("Starting aEye Watcher Core...")
	srv := server.New()
	srv.Start(*port)
}

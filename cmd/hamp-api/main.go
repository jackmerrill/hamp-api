package main

import "github.com/jackmerrill/hamp-api/internal/server"

func main() {
	if err := server.Start(); err != nil {
		panic(err)
	}
}

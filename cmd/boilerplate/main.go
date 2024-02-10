package main

import (
	"log"
	"os"

	"github.com/argatu/boilerplate/internal/api"
)

func main() {
	if err := api.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"

	config "github.com/joao-alho/gator/internal"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read config: %v\n", cfg)

	cfg.SetUser("jalho")

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read config: %v\n", cfg)
}

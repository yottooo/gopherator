package main

import (
	"flag"
	"log"
)

var port = flag.Int("port", 8080, "")

func main() {
	flag.Parse()

	err := initServer(*port)
	if err != nil {
		log.Fatal(err)
	}
}

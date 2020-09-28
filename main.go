package main

import (
	"flag"
	"fmt"
	"gol/world"
	"io"
	"log"
	"os"
)
var input = flag.String("input", "-", "input source to use")
var iteration = flag.Int("generations",1, "number of generations to increment")

func main() {
	var err error
	var reader io.Reader = os.Stdin

	flag.Parse()

	if *input != "-"{
		fh,err := os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()
		reader = fh
	}

	g,err := world.NewWorldFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	for ; *iteration > 0; *iteration-- {
		g = g.NextGeneration()
		fmt.Println(g.String())
	}
}

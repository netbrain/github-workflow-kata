package main

import (
	"flag"
	"fmt"
	"gol/world"
	"io"
	"log"
	"os"
	"time"
)
var input = flag.String("input", "", "input source to use or '-' for stdin")
var iteration = flag.Int("generations",1, "number of generations to increment")
var pdelay = flag.Int("print-delay",0, "delay in ms between print of generations")
var clear = flag.Bool("clear",false, "clear terminal between printout (ansi only)")

func main() {
	var err error
	var reader io.Reader = os.Stdin

	flag.Parse()

	if *input == ""{
		flag.PrintDefaults()
		os.Exit(0)
	}

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
		if *clear {
			//fmt.Print("\033[2J") //clear ansi terminal
			//fmt.Print("\033[H") //put cursor at beginning
			fmt.Print("\033c") //reset terminal
		}
		g = g.NextGeneration()
		fmt.Println(g.String())
		time.Sleep(time.Millisecond*time.Duration(*pdelay))
	}
}

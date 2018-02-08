// TO DO
//  - add text showing iteration

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	N        = 100
	addrFlag = flag.String("addr", ":5555", "server address:port")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	go webServer(addrFlag)

	///////////////////////////////////////////////////////////////
	// Simple example of grid construction and initialization
	grid := NewGrid()
	grid.InitBlabla()
	for i := 0; ; i++ {
		fmt.Println("step", i)
		// 		time.Sleep(1 * time.Millisecond)
		grid.Evolve()
		Plot(grid)
	}
	///////////////////////////////////////////////////////////////
}

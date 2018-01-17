// TO DO
//  - add text showing iteration

package main

import (
	"fmt"
	"math/rand"

	"golang.org/x/exp/shiny/driver"
)

var grid *Grid

func init() {
	grid = NewGrid()
	grid.InitRandom()
}

func main() {
	driver.Main(GridGraph)
}

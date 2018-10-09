package main

import (
	"fmt"

	"github.com/boris-lenzinger/repeatit/tools"
)

func main() {
	s := "1:4"
	a, _ := tools.ParseNumberSerie(s)
	fmt.Printf("%v\n", a)
}

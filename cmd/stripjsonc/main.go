package main

import (
	"os"

	"github.com/dnmfarrell/stripjsonc"
)

func main() {
	stripjsonc.StripJSONCStream(os.Stdin, os.Stdout)
}

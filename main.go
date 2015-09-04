package main

import (
	"os"

	"github.com/rodaine/esu/esu"
)

func main() {
	esu.InitApp().Run(os.Args)
}

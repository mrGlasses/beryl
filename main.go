package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrGlasses/beryl/arguments"
)

func main() {

	result, err := arguments.ExecuteArguments(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(result)
}

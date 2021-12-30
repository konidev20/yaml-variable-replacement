package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magiconair/properties"
	"github.com/rwtodd/Go.Sed/sed"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getProperties(filepath string) *properties.Properties {
	p := properties.MustLoadFile(filepath, properties.UTF8)
	return p
}

func main() {
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)

	inputFilePtr := createCommand.String("f", "", "Input File where variables have to be replaced.")
	varFilePtr := createCommand.String("v", "", "Properties File containing variables.")
	outputFilePtr := createCommand.String("o", "", "Output file where the properties have been replaced. Default = output.txt")
	createCommand.Parse(os.Args[2:])

	outputFile, err := os.OpenFile(*outputFilePtr, os.O_RDWR|os.O_CREATE, 0600)
	check(err)
	defer outputFile.Close()

	inputFile, err := os.Open(*inputFilePtr)
	check(err)
	defer inputFile.Close()

	props := getProperties(*varFilePtr)
	propertyMap := props.Map()

	var inputReader io.Reader = inputFile
	for findVariable, replaceVariable := range propertyMap {
		sedExpression := getSedExpression(findVariable, replaceVariable)
		engine, err := sed.New(strings.NewReader(sedExpression))
		check(err)

		inputReader = engine.Wrap(inputReader)
	}

	n, err := io.Copy(outputFile, inputReader)
	check(err)
	fmt.Println(n)
}

func getSedExpression(findVariable string, replaceVariable string) string {
	sedExpression := fmt.Sprintf("s|%s|%s|g", findVariable, replaceVariable)
	return sedExpression
}

package main

import (
	"fmt"
	"io"
	"jsonParser/parser"
	"log"
	"os"
)

func main() {
	f, err := os.Open("test.json")
	if err != nil {
		panic(err)
	}
	s, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	p := parser.New(s)
	n, err := p.Parse()
	if err != nil {
		panic(err)
	}

	log.Println("Parsing finished")

	output, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	output.WriteString(n.String("\t"))

	fmt.Println(n.Get("education[1].courses[1]", " "))
}

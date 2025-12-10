package main

import (
	"flag"
	"github.com/2540587868/optimus/internal/gen"
	"log"
	"os"
)

func main() {
	// 1. get param
	typePtr := flag.String("type", "", "struct name")
	flag.Parse()

	if *typePtr == "" {
		log.Fatal("Usage: option-gen -type StructName")
	}

	// 2. get env
	goFile := os.Getenv("GOFILE")
	goPackage := os.Getenv("GOPACKAGE")

	// 3. gen
	g := gen.Generator{
		StructName: *typePtr,
		FileName:   goFile,
		Package:    goPackage,
	}

	if err := g.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

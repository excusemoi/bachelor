package main

import (
	"github.com/bachelor/internal/generator"
	"log"
	"path/filepath"
)

func main() {
	var (
		gen = generator.New()
		err error
	)
	if err = gen.Init(filepath.Join("internal", "generator", "configs"), "config-local"); err != nil {
		log.Println("generator")
		log.Fatal(err)
	}
	if err = gen.Run(); err != nil {
		log.Println("generator")
		log.Fatal(err)
	}

}

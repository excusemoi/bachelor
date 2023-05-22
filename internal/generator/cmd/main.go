package main

import (
	"log"
	"path/filepath"

	"github.com/bachelor/internal/generator"
)

func main() {
	var (
		gen = generator.New()
		err error
	)
	if err = gen.Init(filepath.Join("..", "configs"), "config-local"); err != nil {
		log.Println("generator")
		log.Fatal(err)
	}

	if err = gen.Run(); err != nil {
		log.Println("generator")
		log.Fatal(err)
	}

}

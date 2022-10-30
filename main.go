package main

import (
	"flag"
	"log"
	"task-parser/cmd"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "filePath", "", "the path to the task file, the path must be an absolute path")
	flag.Parse()
	if filePath == "" {
		log.Fatal("must provide a path to the task file")
	}
	msg, err := cmd.Run(filePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(msg)
}

package main

import (
	"log"
	"terminal-spinner/spinner"
	"time"
)

func main() {
	s := spinner.New(spinner.Config{})

	log.Println("Spinner started")
	s.Start()

	time.Sleep(time.Second * 10)

	s.Stop()
	log.Println("Spinner stopped")
}
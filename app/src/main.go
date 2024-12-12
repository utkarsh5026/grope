package main

import (
	"fmt"

	"github.com/codecrafters-io/grep-starter-go/src/fw"
)

func main() {
	fmt.Println("Starting file watcher...")
	err := fw.StartWatching(".", func(event fw.FileEvent) error {
		fmt.Println(event.Type.String() + " " + event.Path)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

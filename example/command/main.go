package main

import (
	"fmt"
	"strings"

	"github.com/mittwald/clapper"
)

type Config struct {
	Help    bool   `clapper:"short,long,help='Display help message"`
	Command string `clapper:"command,help=say|sing <message>"`
}

// invoke like `go run ./example/command/main.go -- sing hello world `
func main() {
	var config Config
	trailing, err := clapper.Parse(&config)
	if err != nil {
		config.Help = true
	}

	if err == nil && len(trailing) == 0 {
		config.Help = true
		fmt.Println("Missing command.")
	}

	if config.Help {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("Usage: sample [options] command")
		help, err := clapper.HelpDefault(&config)
		if err != nil {
			panic(err)
		}
		fmt.Println(help)
		return
	}

	message := strings.Join(trailing, " ")
	// do what we are supposed to do with the command
	switch config.Command {
	case "say":
		fmt.Printf("Saying: %s\n", message)
	case "sing":
		fmt.Printf("Singing: %s\n", message)
	default:
		fmt.Printf("Unknown command: %s\n", config.Command)
	}
}

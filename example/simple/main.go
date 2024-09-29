package main

import (
	"fmt"

	"github.com/Dirk007/clapper/pkg/clapper"
)

type Config struct {
	Help    bool    `clapper:"short,long,help='Display help message"`
	Version bool    `clapper:"short,long,help='Display version information'"`
	Debug   bool    `clapper:"short,long,help='Enable debug mode'"`
	Server  string  `clapper:"short,long,default='localhost:8080',help='Server to connect to'"`
	User    string  `clapper:"short,long,help='Username for authentication'"`
	Pass    *string `clapper:"short,long,help='Password will be used for authentication'"`
}

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

	if config.Version {
		fmt.Println("Sample version 1.0.0")
		return
	}

	if config.Debug {
		fmt.Println("Debug mode enabled")
	}

	command := trailing[0]

	fmt.Printf("Parsed config: %+v\n", config)
	fmt.Printf("Command: %s\n", command)
}

package main

import "github.com/airbnb/rudolph/internal/cli"

var version = "development"

func main() {
	cli.Execute(version)
}

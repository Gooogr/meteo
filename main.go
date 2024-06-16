package main

import (
	"meteo/cmd"
	_ "meteo/cmd/set"
)

func main() {
	cmd.Execute()
}

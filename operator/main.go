package main

import "github.com/stackrox/rox/operator/cmd/central-operator/cmd"

func main() {
	//TODO: Based on linker flags switch between sensor or central operator cmd root
	cmd.Execute()
}

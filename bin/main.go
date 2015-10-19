package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"

	"github.com/nlamirault/docker-machine-scaleway"
)

func main() {
	plugin.RegisterDriver(new(scaleway.Driver))
}

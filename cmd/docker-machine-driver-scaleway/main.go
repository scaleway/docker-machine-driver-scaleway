package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	scaleway "github.com/scaleway/docker-machine-driver-scaleway/internal/app/docker-machine-driver-scaleway"
)

func main() {
	plugin.RegisterDriver(scaleway.NewDriver("", ""))
}

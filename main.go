package main

import (
	"github.com/QuentinPerez/docker-machine-driver-scaleway/driver"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(scaleway.NewDriver("", ""))
}

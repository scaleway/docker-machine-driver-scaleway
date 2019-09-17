package dockermachinedriverscaleway

import (
	"errors"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

const (
	// VERSION is the semver version of the package.
	VERSION = "v2.0.0"
)

// Driver implements the docker Driver interface.
// https://godoc.org/github.com/docker/machine/libmachine/drivers#Driver
type Driver struct {
	*drivers.BaseDriver
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{}
}

// NewDriver creates a new Scaleway driver.
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		// FIXME: Add default values here.
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

// DriverName returns the name of the driver.
func (d *Driver) DriverName() string {
	return "scaleway"
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags.
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	return nil
}

// Create configures and starts a scaleway server
func (d *Driver) Create() error {
	return nil
}

// GetSSHHostname returns the IP of the server
func (d *Driver) GetSSHHostname() (string, error) {
	return d.IPAddress, nil
}

// GetState returns the state of the server
func (d *Driver) GetState() (state.State, error) {
	return state.None, nil
}

// GetURL returns IP + docker port
func (d *Driver) GetURL() (string, error) {
	return "", nil
}

// Kill does nothing
func (d *Driver) Kill() error {
	return errors.New("scaleway driver does not support kill")
}

// Remove shutdowns the server and removes the IP
func (d *Driver) Remove() error {
	return nil
}

// Restart reboots the server
func (d *Driver) Restart() error {
	return nil
}

// Start starts the server
func (d *Driver) Start() error {
	return nil
}

// Stop stops the server
func (d *Driver) Stop() error {
	return nil
}

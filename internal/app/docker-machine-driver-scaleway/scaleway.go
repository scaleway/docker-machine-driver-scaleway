package dockermachinedriverscaleway

import (
	"errors"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "SCW_SECRET_KEY",
			Name:   "scaleway-secret-key",
			Usage:  "Scaleway secret key",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_SSH_USER",
			Name:   "scaleway-ssh-user",
			Usage:  "Your SSH user",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_SSH_PORT",
			Name:   "scaleway-ssh-port",
			Usage:  "Your SSH port",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_ORGANIZATION_ID",
			Name:   "scaleway-organization-id",
			Usage:  "Scaleway organzation ID",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_ZONE",
			Name:   "scaleway-zone",
			Usage:  "Scaleway zone (fr-par-1, fr-par-2 or nl-ams-1)",
			Value:  string(scw.ZoneFrPar1),
		},
		mcnflag.BoolFlag{
			EnvVar: "SCW_DEBUG",
			Name:   "scaleway-debug",
			Usage:  "Enable debug mode",
		},

		mcnflag.StringFlag{
			EnvVar: "SCW_NAME",
			Name:   "scaleway-name",
			Usage:  "The name of the server",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_TYPE",
			Name:   "scaleway-type",
			Usage:  "The commercial type of the server",
			Value:  "DEV1-S",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_IMAGE",
			Name:   "scaleway-image",
			Usage:  "The UUID or label of the image for the server",
			Value:  "ubuntu-bionic",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_TAGS",
			Name:   "scaleway-tags",
			Usage:  "Comma-separated list of tags to apply to the server",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_SECURITY_GROUP_ID",
			Name:   "scaleway-security-group-id",
			Usage:  "The security group the server is attached to",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_PLACEMENT_GROUP_ID",
			Name:   "scaleway-placement-group-id",
			Usage:  "The placement group the server is attached to",
		},
		mcnflag.BoolFlag{
			EnvVar: "SCW_ENABLE_IP_V6",
			Name:   "scaleway-enable-ip-v6",
			Usage:  "Determines if IPv6 is enabled for the server",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_IP",
			Name:   "scaleway-ip",
			Usage:  "The flexible IP that you want to attach to your server, UUID or IP value authorized",
		},
	}
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

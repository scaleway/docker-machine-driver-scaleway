package dockermachinedriverscaleway

import (
	"errors"
	"fmt"
	"net"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// TODO:
// - [ ] Configure SDK logger
// - [ ] Create defaults variables

const (
	// VERSION is the semver version of the package.
	VERSION = "v2.0.0"
)

// Driver implements the docker Driver interface.
// https://godoc.org/github.com/docker/machine/libmachine/drivers#Driver
type Driver struct {
	*drivers.BaseDriver
	SecretKey        string
	OrganizationID   string
	Zone             scw.Zone
	Debug            bool
	Name             string
	Type             string
	Image            string
	Tags             []string
	SecurityGroupID  string
	PlacementGroupID string
	EnableIPV6       bool
	IP               string

	instanceUtils *InstanceUtils
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "SCW_SSH_USER",
			Name:   "scaleway-ssh-user",
			Usage:  "Your SSH user",
		},
		mcnflag.IntFlag{
			EnvVar: "SCW_SSH_PORT",
			Name:   "scaleway-ssh-port",
			Usage:  "Your SSH port",
		},
		mcnflag.StringFlag{
			EnvVar: "SCW_SSH_KEY_PATH",
			Name:   "scaleway-ssh-key-path",
			Usage:  "Your SSH private key path",
		},

		mcnflag.StringFlag{
			EnvVar: "SCW_SECRET_KEY",
			Name:   "scaleway-secret-key",
			Usage:  "Scaleway secret key",
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
			EnvVar: "SCW_IP_ID",
			Name:   "scaleway-ip",
			Usage:  "The flexible IP that you want to attach to your server, UUID or IP value authorized", // FIXME: Possible UUID or IP ?
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
	d.SSHUser = flags.String("scaleway-ssh-user")
	d.SSHPort = flags.Int("scaleway-ssh-port")
	d.SSHKeyPath = flags.String("scaleway-ssh-key-path")

	d.SecretKey = flags.String("scaleway-secret-key")
	d.OrganizationID = flags.String("scaleway-organization-id")
	zone, _ := scw.ParseZone(flags.String("scaleway-zone"))
	d.Zone = zone
	d.Debug = flags.Bool("scaleway-debug")

	d.Name = flags.String("scaleway-name")
	d.Type = flags.String("scaleway-type")
	d.Image = flags.String("scaleway-image")
	d.Tags = flags.StringSlice("scaleway-tags")
	d.SecurityGroupID = flags.String("scaleway-security-group-id")
	d.PlacementGroupID = flags.String("scaleway-placement-group-id")
	d.EnableIPV6 = flags.Bool("scaleway-enable-ip-v6")
	d.IP = flags.String("scaleway-ip")

	d.SetSwarmConfigFromFlags(flags)

	return nil
}

// PreCreateCheck is called to enforce pre-creation steps
func (d *Driver) PreCreateCheck() error {
	return nil
}

// TODO: Reomve local host if error
// When their is an error from during the remote host creation
// the local host is neverless created. I can't find any specific
// handeling for this case in other drivers.

// Create configures and starts a scaleway server
func (d *Driver) Create() error {
	log.Info("Generating SSH Key")

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	log.Info("Creating host...")

	instanceUtils, err := d.getInstanceUtils()
	if err != nil {
		return err
	}

	err = instanceUtils.CreateServer()
	if err != nil {
		return fmt.Errorf("cannot create the host: %s", err)
	}

	return nil
}

// GetSSHHostname returns the IP of the server
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetState returns the state of the server
func (d *Driver) GetState() (state.State, error) {
	instanceUtils, err := d.getInstanceUtils()
	if err != nil {
		return state.None, err
	}

	server, err := instanceUtils.GetCreatedServer()
	if err != nil {
		return state.Error, err
	}

	switch server.State {
	case instance.ServerStateRunning:
		return state.Running, nil
	case instance.ServerStateStopped, instance.ServerStateStoppedInPlace:
		return state.Stopped, nil
	case instance.ServerStateStarting:
		return state.Starting, nil
	case instance.ServerStateStopping:
		return state.Stopping, nil
	}

	return state.Error, nil
}

// GetURL returns IP + docker port
func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// Kill does nothing
func (d *Driver) Kill() error {
	return errors.New("scaleway driver does not support kill")
}

// Remove shutdowns the server and removes the IP
func (d *Driver) Remove() error {
	instanceUtils, err := d.getInstanceUtils()
	if err != nil {
		return err
	}

	return instanceUtils.RemoveServer()
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

func (d *Driver) getInstanceUtils() (*InstanceUtils, error) {
	if d.instanceUtils == nil {
		instanceUtils, err := NewInstanceUtils(d)
		if err != nil {
			return nil, err
		}

		d.instanceUtils = instanceUtils
	}

	return d.instanceUtils, nil
}

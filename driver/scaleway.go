package scaleway

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
)

const (
	// VERSION represents the semver version of the package
	VERSION           = "v1.0.0+dev"
	defaultImage      = "ubuntu-trusty"
	defaultBootscript = "docker"
)

var scwAPI *api.ScalewayAPI

// Driver represents the docker driver interface
type Driver struct {
	*drivers.BaseDriver
	ServerID       string
	Organization   string
	IPID           string
	Token          string
	commercialType string
	name           string
	stopping       bool
	// size         string
	// userDataFile string
	// ipv6         bool
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "scaleway"
}

func (d *Driver) getClient() (cl *api.ScalewayAPI, err error) {
	if scwAPI == nil {
		scwAPI, err = api.NewScalewayAPI(d.Organization, d.Token, "docker-machine-driver-scaleway/%v"+VERSION)
	}
	cl = scwAPI
	return
}

// SetConfigFromFlags sets the flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) (err error) {
	d.Token, d.Organization = flags.String("scaleway-token"), flags.String("scaleway-organization")
	if d.Token == "" || d.Organization == "" {
		config, cfgErr := config.GetConfig()
		if cfgErr == nil {
			if d.Token == "" {
				d.Token = config.Token
			}
			if d.Organization == "" {
				d.Organization = config.Organization
			}
		} else {
			return fmt.Errorf("You must provide organization and token")
		}
	}
	d.commercialType = flags.String("scaleway-commercial-type")
	d.name = flags.String("scaleway-name")
	return
}

// NewDriver returns a new driver
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{},
	}
}

// GetCreateFlags registers the flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_TOKEN",
			Name:   "scaleway-token",
			Usage:  "Scaleway token",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_ORGANIZATION",
			Name:   "scaleway-organization",
			Usage:  "Scaleway organization",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_NAME",
			Name:   "scaleway-name",
			Usage:  "Assign a name",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_COMMERCIAL_TYPE",
			Name:   "scaleway-commercial-type",
			Usage:  "Specifies the commercial type",
			Value:  "VC1S",
		},
		// mcnflag.StringFlag{
		//     EnvVar: "SCALEWAY_USERDATA",
		//     Name:   "scaleway-userdata",
		//     Usage:  "Path to file with user-data",
		// },
		// mcnflag.BoolFlag{
		// 	EnvVar: "SCALEWAY_IPV6",
		// 	Name:   "scaleway-ipv6",
		// 	Usage:  "Enable ipv6",
		// },
	}
}

// Create configures and starts a scaleway server
func (d *Driver) Create() (err error) {
	var publicKey []byte
	var cl *api.ScalewayAPI
	var ip *api.ScalewayGetIP

	log.Infof("Creating SSH key...")
	if err = ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}
	publicKey, err = ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return
	}
	log.Infof("Creating server...")
	cl, err = d.getClient()
	if err != nil {
		return
	}
	ip, err = cl.NewIP()
	if err != nil {
		return
	}
	d.IPAddress = ip.IP.Address
	d.ServerID, err = api.CreateServer(cl, &api.ConfigCreateServer{
		ImageName:      defaultImage,
		CommercialType: d.commercialType,
		Name:           d.name,
		Bootscript:     defaultBootscript,
		IP:             ip.IP.ID,
		Env: strings.Join([]string{"AUTHORIZED_KEY",
			strings.Replace(string(publicKey[:len(publicKey)-1]), " ", "_", -1)}, "="),
	})
	if err != nil {
		return
	}
	log.Infof("Starting server...")
	err = api.StartServer(cl, d.ServerID, false)
	return
}

// GetSSHHostname returns the IP of the server
func (d *Driver) GetSSHHostname() (string, error) {
	return d.IPAddress, nil
}

// GetState returns the state of the server
func (d *Driver) GetState() (st state.State, err error) {
	var server *api.ScalewayServer
	var cl *api.ScalewayAPI

	st = state.Error
	cl, err = d.getClient()
	if err != nil {
		return
	}
	server, err = cl.GetServer(d.ServerID)
	if err != nil {
		return
	}
	st = state.None
	switch server.State {
	case "starting":
		st = state.Starting
	case "running":
		st = state.Running
	case "stopping":
		st = state.Stopping
	case "stopped":
		st = state.Stopped
	}
	if d.stopping {
		time.Sleep(5 * time.Second)
	}
	return
}

// GetURL returns IP + docker port
func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(d.IPAddress, "2376")), nil
}

func (d *Driver) postAction(action string) (err error) {
	var cl *api.ScalewayAPI

	cl, err = d.getClient()
	if err != nil {
		return
	}
	err = cl.PostServerAction(d.ServerID, action)
	return
}

// Kill does nothing
func (d *Driver) Kill() error {
	return errors.New("scaleway driver does not support kill")
}

// Remove shutdowns the server and removes the IP
func (d *Driver) Remove() (err error) {
	var cl *api.ScalewayAPI

	cl, err = d.getClient()
	if err != nil {
		return
	}
	err = cl.PostServerAction(d.ServerID, "terminate")
	if err != nil {
		return
	}
	for {
		_, err = cl.GetServer(d.ServerID)
		if err != nil {
			break
		}
	}
	err = cl.DeleteIP(d.IPID)
	return
}

// Restart reboots the server
func (d *Driver) Restart() error {
	return d.postAction("reboot")
}

// Start starts the server
func (d *Driver) Start() error {
	return d.postAction("poweron")
}

// Stop stops the server
func (d *Driver) Stop() error {
	d.stopping = true
	return d.postAction("poweroff")
}

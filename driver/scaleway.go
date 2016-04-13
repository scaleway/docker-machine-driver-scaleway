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
)

const (
	// VERSION represents the semver version of the package
	VERSION           = "v1.0.0"
	defaultImage      = "ubuntu-trusty"
	defaultBootscript = "docker"
)

var scwAPI *api.ScalewayAPI

type Driver struct {
	*drivers.BaseDriver
	ServerID       string
	Organization   string
	Token          string
	commercialType string
	name           string
	// size         string
	// userDataFile string
	// ipv6         bool
}

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

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) (err error) {
	d.Token, d.Organization = flags.String("scaleway-token"), flags.String("scaleway-organization")
	if d.Token == "" || d.Organization == "" {
		return fmt.Errorf("You must provide organization and token")
	}
	d.commercialType = flags.String("scaleway-commercial-type")
	d.name = flags.String("scaleway-name")
	return
}

func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{},
	}
}

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

func (d *Driver) GetSSHHostname() (string, error) {
	return d.IPAddress, nil
}

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
	time.Sleep(5 * time.Second)
	return
}

func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(d.IPAddress, "2376")), nil
}

func (d *Driver) Kill() error {
	return errors.New("scaleway driver does not support kill")
}

func (d *Driver) Remove() error {
	log.Info("Remove: not implemented yet")
	return nil
}

func (d *Driver) Restart() error {
	log.Info("Restart: not implemented yet")
	return nil
}

func (d *Driver) Start() error {
	log.Info("Start: not implemented yet")
	return nil
}

func (d *Driver) Stop() error {
	log.Info("Stop: not implemented yet")
	return nil
}

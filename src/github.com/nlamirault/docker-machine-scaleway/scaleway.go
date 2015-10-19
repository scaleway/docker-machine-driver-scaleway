// Copyright (C) 2015  Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scaleway

import (
	//"errors"
	"fmt"
	//"os"
	// "os/exec"
	"path/filepath"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

const (
	driverName      = "scaleway"
	version         = "0.1.0"
	dockerConfigDir = "/etc/docker"
)

// Driver defines how a host is created and controlled
// See github.com/docker/machine/libmachine/drivers
type Driver struct {
	*drivers.BaseDriver
	ID           string
	UserID       string
	Token        string
	Organization string
	Image        string
	Volumes      string
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.UserID = flags.String("scaleway-user-id")
	d.Token = flags.String("scaleway-token")
	d.Organization = flags.String("scaleway-organization")
	d.Image = flags.String("scaleway-image")
	d.Volumes = flags.String("scaleway-volumes")
	d.SSHUser = "docker"
	d.SSHPort = 22
	if d.UserID == "" {
		return fmt.Errorf("scaleway driver requires the --scaleway-userid option")
	}
	if d.Token == "" {
		return fmt.Errorf("scaleway driver requires the --scaleway-token option")
	}
	if d.Organization == "" {
		return fmt.Errorf("scaleway driver requires the --scaleway-organization option")
	}
	if d.Image == "" {
		return fmt.Errorf("scaleway driver requires the --scaleway-image option")
	}
	return nil
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.Flag{
			EnvVar: "SCALEWAY_USER_ID",
			Name:   "scaleway-user-id",
			Usage:  "ID of the Scaleway user",
			Value:  "",
		},
		mcnflag.Flag{
			EnvVar: "SCALEWAY_TOKEN",
			Name:   "scaleway-token",
			Usage:  "Token for the Scaleway API",
			Value:  "",
		},
		mcnflag.Flag{
			EnvVar: "SCALEWAY_ORGANIZATION",
			Name:   "scaleway-organization",
			Usage:  "Scaleway Organization",
			Value:  "",
		},
		mcnflag.Flag{
			EnvVar: "SCALEWAY_IMAGE",
			Name:   "scaleway-image",
			Usage:  "Scaleway image to use for machine",
			Value:  "",
		},
	}
}

func (d *Driver) DriverName() string {
	return driverName
}

func (d *Driver) GetMachineName() string {
	return d.MachineName
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHKeyPath() string {
	// return filepath.Join(d.LocalArtifactPath("."), "id_rsa")
	return filepath.Join(".", "id_rsa")
}

func (d *Driver) GetSSHPort() (int, error) {
	if d.SSHPort == 0 {
		d.SSHPort = 22
	}

	return d.SSHPort, nil
}

func (d *Driver) GetSSHUsername() string {
	if d.SSHUser == "" {
		d.SSHUser = "docker"
	}

	return d.SSHUser
}

func (d *Driver) GetState() (state.State, error) {
	log.Debugf("[Scaleway] Retrieving state server %s", d.ID)
	client := d.getClient()
	response, err := client.GetServer(d.ID)
	if err != nil {
		return state.Error, err
	}
	return getServerState(response.Server.State), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) PreCreateCheck() error {
	// Others...?
	return nil
}

func (d *Driver) Create() error {
	log.Infof("[Scaleway] Creating instance...")
	client := d.getClient()
	response, err := client.CreateServer(
		d.MachineName, d.Image)
	if err != nil {
		return err
	}
	d.ID = response.Server.ID
	log.Debugf("[Scaleway] ServerID %s", d.ID)

	log.Debugf("[Scaleway] Create SSH key")
	err = d.createSSHKey()
	if err != nil {
		return err
	}
	log.Debugf("[Scaleway] Upload SSH key")
	if _, err = client.UploadPublicKey(d.publicSSHKeyPath()); err != nil {
		return err
	}

	if err = d.Start(); err != nil {
		return err
	}
	log.Debugf("[Scaleway] Waiting server ready .......")
	if err := d.waitForServerState(state.Running); err != nil {
		return err
	}
	// log.Debugf("[Scaleway] Waiting SSH .......")
	// if err := ssh.WaitForTCP(fmt.Sprintf("%s:%d", d.IPAddress, 22)); err != nil {
	// 	return err
	// }

	time.Sleep(10 * time.Second)
	d.setupHostname()
	//d.installDocker()
	return nil
}

func (d *Driver) Start() error {
	log.Infof("[Scaleway] Starting instance...")
	client := d.getClient()
	if _, err := client.PerformServerAction(d.ID, "poweron"); err != nil {
		return err
	}
	d.waitForServerState(state.Running)
	return nil
}

func (d *Driver) Stop() error {
	log.Infof("[Scaleway] Stopping instance...")
	client := d.getClient()
	if _, err := client.PerformServerAction(d.ID, "poweroff"); err != nil {
		return err
	}
	d.waitForServerState(state.Stopped)
	return nil
}

func (d *Driver) Remove() error {
	log.Infof("[Scaleway] Removing instance... ")
	client := d.getClient()
	if err := client.DeleteServer(d.ID); err != nil {
		return err
	}
	d.waitForServerState(state.Stopped)
	return nil
}

func (d *Driver) Restart() error {
	log.Infof("[Scaleway] Rebooting instance...")
	client := d.getClient()
	if _, err := client.PerformServerAction(d.ID, "reboot"); err != nil {
		return err
	}
	d.waitForServerState(state.Running)
	return nil
}

func (d *Driver) Kill() error {
	return d.Stop()
}

func (d *Driver) setupHostname() error {
	log.Debugf("[Scaleway] Setting hostname: %s", d.MachineName)
	_, err := drivers.RunSSHCommandFromDriver(d,
		fmt.Sprintf(
			"echo \"127.0.0.1 %s\" | sudo tee -a /etc/hosts && sudo hostname %s && echo \"%s\" | sudo tee /etc/hostname",
			d.MachineName,
			d.MachineName,
			d.MachineName,
		))
	return err
}

func (d *Driver) createSSHKey() error {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}
	return nil
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) getClient() *ScalewayClient {
	return NewClient(d.Token, d.UserID, d.Organization)

}

func getServerState(status string) state.State {
	switch status {
	case "stopped":
		return state.Stopped
	case "stopping":
		return state.Stopping
	case "starting":
		return state.Starting
	case "running":
		return state.Running
	}
	return state.None
}

func (d *Driver) waitForServerState(serverState state.State) error {
	client := d.getClient()
	for {
		response, err := client.GetServer(d.ID)
		if err != nil {
			return err
		}
		status := getServerState(response.Server.State)
		log.Infof("[Scaleway] Waiting server state %s. Currently : %s",
			serverState, status)
		if status == serverState {
			if status == state.Running {
				log.Infof("[Scaleway] Server %s is running. Set IP address",
					d.ID)
				d.IPAddress = response.Server.PublicIP.Address
			}
			break
		}
		time.Sleep(5 * time.Second)
	}
	return nil

}

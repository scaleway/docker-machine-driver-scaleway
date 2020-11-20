package scaleway

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"

	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	"github.com/sirupsen/logrus"
)

const (
	// VERSION represents the semver version of the package
	VERSION      = "v2.0.0"
	defaultImage = "ubuntu-focal"
	DELAY        = 30 // in second
)

var instanceApi *instance.API

// Driver represents the docker driver interface
type Driver struct {
	*drivers.BaseDriver
	ServerID       string
	OrganizationID string
	ProjectID      string
	IPID           string
	AccessKey      string
	SecretKey      string
	CommercialType string
	Zone           scw.Zone
	name           string
	image          string
	bootscript     string
	ip             string
	volumes        string
	IPPersistant   bool
	stopping       bool
	created        bool
	ipv6           bool
	// userDataFile string
}

// CloudConfigUser represents the user of "users" cloud-init config
type CloudConfigUser struct {
	Name              string   `yaml:"name,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	LockPasswd        bool     `yaml:"lock_passwd,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Sudo              string   `yaml:"sudo,omitempty"`
}

// CloudConfigUsers represents "users" cloud-init config
type CloudConfigUsers struct {
	Users []CloudConfigUser `yaml:"users,omitempty"`
}

// CloudConfigPackages represents "packages" cloud-init config
type CloudConfigPackages struct {
	Packages []string `yaml:"packages,omitempty"`
}

// type CloudConfigRunCmd struct {
// 	RunCmd []CloudConfigCmd `yaml:"runcmd,omitempty"`
// }

// type CloudConfigCmd struct {
// 	Cmd string `yaml:",flow"`
// }

// NewDriver returns a new driver
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{},
	}
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	if d.CommercialType == "" {
		return "scaleway-v2"
	}
	return fmt.Sprintf("scaleway(%v)", d.CommercialType)
}

// GetCreateFlags registers the flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_ACCESSKEY",
			Name:   "scaleway-accesskey",
			Usage:  "Scaleway accesskey (required)",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_SECREYKEY",
			Name:   "scaleway-secretkey",
			Usage:  "Scaleway secretkey (required)",
		},
		// mcnflag.StringFlag{
		// 	EnvVar: "SCALEWAY_ORGANIZATION_ID",
		// 	Name:   "scaleway-organization-id",
		// 	Usage:  "Scaleway organization id",
		// },
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_PROJECT_ID",
			Name:   "scaleway-project-id",
			Usage:  "Scaleway project id (required)",
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
			Value:  "DEV1-S",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_ZONE",
			Name:   "scaleway-zone",
			Usage:  "Specifies the location (fr-par-1, fr-par-2, nl-ams-1, pl-waw-1)",
			Value:  "fr-par-2",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_IMAGE",
			Name:   "scaleway-image",
			Usage:  "Specifies the image",
			Value:  defaultImage,
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_BOOTSCRIPT",
			Name:   "scaleway-bootscript",
			Usage:  "Specifies the bootscript",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "SCALEWAY_IP",
			Name:   "scaleway-ip",
			Usage:  "Specifies the Public IP address",
			Value:  "",
		},
		// mcnflag.StringFlag{
		// 	EnvVar: "SCALEWAY_VOLUMES",
		// 	Name:   "scaleway-volumes",
		// 	Usage:  "Attach additional volume (e.g., 50G)",
		// 	Value:  "",
		// },
		// mcnflag.StringFlag{
		// 	EnvVar: "SCALEWAY_USER",
		// 	Name:   "scaleway-user",
		// 	Usage:  "Specifies SSH user name",
		// 	Value:  drivers.DefaultSSHUser,
		// },
		// mcnflag.IntFlag{
		// 	EnvVar: "SCALEWAY_PORT",
		// 	Name:   "scaleway-port",
		// 	Usage:  "Specifies SSH port",
		// 	Value:  drivers.DefaultSSHPort,
		// },
		mcnflag.BoolFlag{
			EnvVar: "SCALEWAY_DEBUG",
			Name:   "scaleway-debug",
			Usage:  "Enables Scaleway client debugging",
		},
		mcnflag.BoolFlag{
			EnvVar: "SCALEWAY_IPV6",
			Name:   "scaleway-ipv6",
			Usage:  "Enable ipv6",
		},
		// mcnflag.StringFlag{
		//     EnvVar: "SCALEWAY_USERDATA",
		//     Name:   "scaleway-userdata",
		//     Usage:  "Path to file with user-data",
		// },
	}
}

// SetConfigFromFlags sets the flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) (err error) {
	// log.Infof("SetConfigFromFlags ...")
	if flags.Bool("scaleway-debug") {
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(logrus.DebugLevel)
		log.SetDebug(true)
		log.Info("Set Log to DEBUG")
		log.Debug("debug log: on")
	}

	d.AccessKey = flags.String("scaleway-accesskey")
	d.SecretKey = flags.String("scaleway-secretkey")
	if d.AccessKey == "" || d.SecretKey == "" {
		return fmt.Errorf("You must provide accesskey and secretkey")
	}

	/*
	 Prefer use Project instead organization id
	*/
	d.ProjectID = flags.String("scaleway-project-id")
	if d.ProjectID == "" {
		return fmt.Errorf("You must provide project-id or organization-id")
	}

	// d.OrganizationID = flags.String("scaleway-organization-id")
	// if d.ProjectID == "" && d.OrganizationID == "" {
	// 	return fmt.Errorf("You must provide project-id or organization-id")
	// }

	zone := flags.String("scaleway-zone")
	d.Zone, err = scw.ParseZone(zone)
	if err != nil {
		return fmt.Errorf("You must provide an known zone")
	}

	d.CommercialType = flags.String("scaleway-commercial-type")
	d.name = flags.String("scaleway-name")
	d.image = flags.String("scaleway-image")
	d.bootscript = flags.String("scaleway-bootscript")
	d.ip = flags.String("scaleway-ip")
	// d.volumes = flags.String("scaleway-volumes")
	d.ipv6 = flags.Bool("scaleway-ipv6")
	// d.BaseDriver.SSHUser = flags.String("scaleway-user")
	// d.BaseDriver.SSHPort = flags.Int("scaleway-port")
	d.BaseDriver.SSHUser = drivers.DefaultSSHUser
	d.BaseDriver.SSHPort = drivers.DefaultSSHPort
	return
}

// Create configures and starts a scaleway server
func (d *Driver) Create() (err error) {
	log.Infof("Create Scaleway Server ...")

	d.getClient()
	log.Infof("Creating SSH key...")
	if err = ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	log.Debugf("Creating config server...")
	bootType := instance.BootTypeLocal
	bootscript := ""
	ipv6 := false
	// var tags []string
	// securityGroup := ""
	// placementGroup := ""
	projectId := &d.ProjectID
	// var organizationId *string = &d.OrganizationID
	var commercialType string = d.CommercialType
	var name string = d.name
	var image string = d.image
	var zone scw.Zone = d.Zone
	log.Debugf("Creating server bootType: %s", bootType)
	log.Debugf("Creating server Bootscript: %s", bootscript)
	log.Debugf("Creating server projectId: %s", string(*projectId))
	// log.Infof("Creating server organizationId: %s", string(*organizationId))
	log.Debugf("Creating server commercialType: %s", commercialType)
	log.Debugf("Creating server name: %s", name)
	log.Debugf("Creating server image: %s", image)
	log.Debugf("Creating server zone: %s", zone)

	config := instance.CreateServerRequest{
		Zone: zone,
		Name: name,
		// DynamicIPRequired: &dynamicIP,
		CommercialType: commercialType,
		Image:          image,
		// Volumes:          		volumes,
		EnableIPv6: ipv6,
		BootType:   &bootType,
		// Bootscript: &bootscript,
		// Organization:   organizationId,
		Project: projectId,
		// Tags:           tags,
		// SecurityGroup:  &securityGroup,
		// PlacementGroup: &placementGroup,

	}
	log.Debugf("Config server created")
	if d.bootscript != "" {
		bootType = instance.BootTypeBootscript
		config.BootType = &bootType
		config.Bootscript = &d.bootscript
	}

	ipRequired := true
	d.IPPersistant = !ipRequired
	config.DynamicIPRequired = &ipRequired

	if d.ip != "" {
		log.Debugf("public ip from conf...")
		if err = d.resolvePublicIP(); err == nil {
			ipRequired = false
			d.IPPersistant = !ipRequired
			config.DynamicIPRequired = &ipRequired
			config.PublicIP = &d.IPID
		}
	}

	log.Debugf("API call: CreateServer %v", config)
	var server *instance.CreateServerResponse
	server, err = instanceApi.CreateServer(&config)
	if err != nil {
		log.Errorf("Server creation error: %s", err)
		return
	}
	log.Debugf("Server created: ", server.Server)
	log.Infof("Server created: %s (%s)", server.Server.Name, server.Server.ID)
	d.ServerID = server.Server.ID

	if d.IPID == "" {
		_ = d.createPublicIP()
	}

	log.Debug("Setting cloud-init config... ")
	userCloudInit, _ := d.cloudInit()
	log.Debugf("userCloudInit:\n%s", string(userCloudInit))
	userDataRequest := instance.SetServerUserDataRequest{
		Zone:     d.Zone,
		ServerID: server.Server.ID,
		Key:      "cloud-init",
		Content:  strings.NewReader(userCloudInit),
	}

	log.Debugf("API call: SetServerUserData %v", userDataRequest)
	err = instanceApi.SetServerUserData(&userDataRequest)
	if err != nil {
		return
	}

	d.Start()
	// log.Infof("Starting server...")
	// var serverAction instance.ServerAction = "poweron" // default action
	// var serverActionRequest = instance.ServerActionRequest{
	// 	Zone:     d.Zone,
	// 	ServerID: d.ServerID,
	// 	Action:   serverAction,
	// }

	// log.Debugf("API call: SetServerUserData %v", userDataRequest)
	// _, err = instanceApi.ServerAction(&serverActionRequest)
	// if err != nil {
	// 	d.created = true
	// }

	return
}

// GetSSHHostname returns the IP of the server
func (d *Driver) GetSSHHostname() (string, error) {
	log.Debugf("GetSSHHostname... %s", d.IPAddress)
	return d.IPAddress, nil
}

// GetState returns the state of the server
func (d *Driver) GetState() (st state.State, err error) {
	log.Debugf("GetState ...")
	d.getClient()

	var serverResponse *instance.GetServerResponse
	var serverRequest = instance.GetServerRequest{
		Zone:     d.Zone,
		ServerID: d.ServerID,
	}
	log.Debugf("API call: GetServer %v", serverRequest)
	serverResponse, err = instanceApi.GetServer(&serverRequest)

	if err != nil {
		d.created = false
		st = state.Error
		return
	}
	log.Debugf("Server state: %s", serverResponse.Server.State.String())

	st = state.None
	d.stopping = false
	d.created = true
	switch serverResponse.Server.State {
	case "starting":
		st = state.Starting
		log.Debugf("Delay %d", DELAY)
		time.Sleep(DELAY * time.Second)
	case "running":
		st = state.Running
	case "stopping":
		st = state.Stopping
		d.stopping = true
		log.Debugf("Delay %d", DELAY)
		time.Sleep(DELAY * time.Second)
	case "stopped":
		st = state.Stopped
	default:
		st = state.Error
		d.created = false
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

// getClient return Scaleway Instance Client
func (d *Driver) getClient() (err error) {
	log.Debug("getClient ...")
	if instanceApi == nil {
		log.Debug("getClient create client ...")
		// region, _ := d.Zone.Region()
		client, errClient := scw.NewClient(
			scw.WithAuth(d.AccessKey, d.SecretKey),
			// scw.WithDefaultRegion(region),
			scw.WithDefaultZone(d.Zone),
		)
		log.Debugf("getClient created client %s", client)
		instanceApi = instance.NewAPI(client)
		err = errClient
	}
	log.Debugf("getClient instanceApi: %s", instanceApi)
	return
}

/*
	cloudInit used to add docker-machine sshkey to Scaleway Instance

	Multiple choices are possible:
		- Use Modules "users":
			- without runcmd to reload cloud-init config no sshkey added to root user :thinking: why ? so KO
			- with runcmd to reload cloud-init config to apply config, sshkey is added to root user but sshkey from scaleway are setted to 'no-port-forwarding, ...'  :thinking: why ? so OK and KO
		- Use Scaleway implementation with instance_keys file and scw-fetch-ssh-keys command is OK

	***** It's very "sensible"...  *****
*/
func (d *Driver) cloudInit() (contentByte string, err error) {
	pub := d.GetSSHKeyPath() + ".pub"
	publicKey, err := ioutil.ReadFile(pub)
	if err != nil {
		return
	}
	log.Debugf("SSH key pub: %s --", string(publicKey))

	s := strings.TrimSpace(string(publicKey))
	// var sshkey []string
	// sshkey = append(sshkey, s)
	// cloudSSHUser := CloudConfigUser{
	// 	Name:              "root",
	// 	SshAuthorizedKeys: sshkey,
	// 	// Sudo:              "ALL=(ALL) NOPASSWD:ALL",
	// 	// Groups:            "sudo",
	// }
	// var users []CloudConfigUser

	// cInit := CloudConfigUsers{
	// 	Users: append(users, cloudSSHUser),
	// }
	// contentUsers, _ := yaml.Marshal(&cInit)

	// var packages []string
	// packages = append(packages, "sudo")
	// cloudConfigPackages := CloudConfigPackages{
	// 	Packages: packages,
	// }
	// contentPackages, _ := yaml.Marshal(&cloudConfigPackages)
	// 	contentByte = fmt.Sprintf(`#cloud-config
	// final_message: "Scaleway is happy to welcome you to a cloud-init enabled instance"
	// runcmd:
	// - [ cloud-init, clean ]
	// - [ cloud-init, init, '--local' ]
	// - [ cloud-init, init ]
	// - [ cloud-init, status ]
	// package_update: true
	// %s
	// %s
	// `,
	// 		contentPackages,
	// 		contentUsers)
	// 	contentByte = fmt.Sprintf(`#cloud-config
	// final_message: "Scaleway is happy to welcome you to a cloud-init enabled instance"
	// runcmd:
	// - [ cloud-init, clean ]
	// - [ cloud-init, init, '--local' ]
	// - [ cloud-init, init ]
	// - [ cloud-init, status ]
	// ssh_authorized_keys:
	// - %s
	// `, s)
	contentByte = fmt.Sprintf(`#cloud-config
final_message: "Scaleway is happy to welcome you to a cloud-init enabled instance"
package_update: true
runcmd:
- [scw-fetch-ssh-keys, '--upgrade']
write_files:
- content: |
    %s
  path: /root/.ssh/instance_keys
  append: true
`,
		s)
	return
}

// resolvePublicIP check and use an Public IP from params
func (d *Driver) resolvePublicIP() (err error) {
	log.Debugf("Get Public IP with conf: %s", d.ip)
	var ips *instance.GetIPResponse

	d.IPPersistant = true
	getIPRequest := instance.GetIPRequest{
		Zone: d.Zone,
		IP:   d.ip,
	}
	log.Debugf("API call: GetIP %v", getIPRequest)
	ips, err = instanceApi.GetIP(&getIPRequest)
	if err != nil {
		d.IPPersistant = false
		log.Errorf("Get Public IP: %s", err.Error())
		return
	}
	if ips.IP.Address.String() == d.ip {
		d.IPAddress = ips.IP.Address.String()
		d.IPID = ips.IP.ID
	}

	return
}

// createPublicIP create and use an Public IP
func (d *Driver) createPublicIP() (err error) {
	log.Debug("Create Public IP")
	var ip *instance.CreateIPResponse
	d.IPPersistant = false

	createIPRequest := instance.CreateIPRequest{
		Zone:    d.Zone,
		Project: &d.ProjectID,
		Server:  &d.ServerID,
	}
	log.Debugf("API call: CreateIP %v", createIPRequest)
	ip, err = instanceApi.CreateIP(&createIPRequest)
	if err != nil {
		log.Errorf("Create Public IP: %s", err.Error())
		return
	}
	d.IPAddress = ip.IP.Address.String()
	d.IPID = ip.IP.ID

	return
}

// postAction post Action to Scaleway Instance API
func (d *Driver) postAction(action instance.ServerAction) (err error) {
	log.Debugf("postAction %s ...", action)
	d.getClient()
	var serverActionRequest = instance.ServerActionRequest{
		Zone:     d.Zone,
		ServerID: d.ServerID,
		Action:   action,
	}

	log.Debugf("API call: ServerAction %v", serverActionRequest)
	_, err = instanceApi.ServerAction(&serverActionRequest)

	return
}

// Kill does nothing
func (d *Driver) Kill() error {
	return errors.New("scaleway driver does not support kill")
}

// Remove shutdowns the server and removes the IP
func (d *Driver) Remove() (err error) {
	d.getClient()
	serverStatus, _ := d.GetState()

	for serverStatus != state.Stopped && serverStatus != state.Error {
		if serverStatus != state.Stopping && serverStatus != state.Stopped {
			d.Stop()
		}
		serverStatus, _ = d.GetState()
	}

	var serverResponse *instance.GetServerResponse
	var serverRequest = instance.GetServerRequest{
		Zone:     d.Zone,
		ServerID: d.ServerID,
	}
	log.Debugf("API call: GetServer %v", serverRequest)
	serverResponse, err = instanceApi.GetServer(&serverRequest)

	if err != nil {
		return
	}

	var deleteServerRequest = instance.DeleteServerRequest{
		Zone:     d.Zone,
		ServerID: d.ServerID,
	}
	log.Infof("Delete server: %s", d.ServerID)
	log.Debugf("API call: GetServer %v", deleteServerRequest)
	errRemove := instanceApi.DeleteServer(&deleteServerRequest)

	volumes := serverResponse.Server.Volumes
	log.Debugf("Delete volumes: %s", volumes)
	for _, volume := range volumes {
		log.Infof("Delete volume: %s", volume.ID)
		deleteVolumeRequest := instance.DeleteVolumeRequest{
			Zone:     volume.Zone,
			VolumeID: volume.ID,
		}
		log.Debugf("API call: DeleteVolume %v", deleteVolumeRequest)
		err = instanceApi.DeleteVolume(&deleteVolumeRequest)
	}

	log.Debugf("Persist IP %s %s %s %s", d.IPPersistant, serverResponse.Server.PublicIP.ID, serverResponse.Server.PublicIP.Address.String(), serverResponse.Server.DynamicIPRequired)
	// if !d.IPPersistant {
	if serverResponse.Server.DynamicIPRequired {
		log.Infof("Delete public IP: %s (%s)", d.IPAddress, d.IPID)
		// log.Infof("Delete IP: %s (%s) ? %s", d.IPAddress, d.IPID, serverResponse.Server.PublicIP.ID)
		deleteIPRequest := instance.DeleteIPRequest{
			Zone: d.Zone,
			IP:   d.IPID,
		}
		log.Debugf("API call: DeleteIP %v", deleteIPRequest)
		err = instanceApi.DeleteIP(&deleteIPRequest)
		if err != nil {
			log.Errorf("Delete IP %s: %s", d.IPAddress, err)
		}
	}
	if errRemove != nil {
		log.Errorf("Delete server %s: %s", d.ServerID, err.Error())
		err = errRemove
	}
	return
}

// Restart reboots the server
func (d *Driver) Restart() error {
	log.Info("Restart Server ...")
	return d.postAction(instance.ServerActionReboot)
}

// Start starts the server
func (d *Driver) Start() error {
	log.Info("Start Server ...")
	return d.postAction(instance.ServerActionPoweron)
}

// Stop stops the server
func (d *Driver) Stop() error {
	log.Info("Stop Server ...")
	d.stopping = true
	return d.postAction(instance.ServerActionPoweroff)
}

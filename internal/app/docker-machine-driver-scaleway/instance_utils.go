package dockermachinedriverscaleway

import (
	"fmt"

	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/version"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/namegenerator"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

type InstanceUtils struct {
	client      *scw.Client
	instanceAPI *instance.API
	driver      *Driver
	serverID    string
}

func NewInstanceUtils(d *Driver) (*InstanceUtils, error) {
	i := &InstanceUtils{
		driver: d,
	}

	log.Debug("Try to migrate config")
	_, err := scw.MigrateLegacyConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot migrate configuration: %s", err)
	}

	log.Debug("Creating Scaleway client")
	config, err := scw.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot load SDK config: %s", err)
	}

	profile, err := config.GetActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("cannot get SDK config active profile: %s", err)
	}

	clientOptions := []scw.ClientOption{
		scw.WithEnv(),
		scw.WithProfile(profile),
		scw.WithUserAgent(fmt.Sprintf("docker-machine/%s", version.Version)),
	}

	if d.Zone != "" {
		clientOptions = append(clientOptions, scw.WithDefaultZone(d.Zone))
	}

	client, err := scw.NewClient(clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("cannot create an SDK client: %s", err)
	}

	i.client = client
	if err != nil {
		return nil, err
	}

	log.Debug("Creating instance API")
	i.instanceAPI = instance.NewAPI(i.client)

	return i, nil
}

func (i *InstanceUtils) getIPIDFromAddress(ipAddress string) (string, error) {
	log.Infof("Finding IP ID for %s", ipAddress)

	res, err := i.instanceAPI.ListIPs(&instance.ListIPsRequest{})
	if err != nil {
		return "", fmt.Errorf("cannot find IP address ID: %s", err)
	}

	for _, IP := range res.IPs {
		if IP.Address.String() == ipAddress {
			return IP.ID, nil
		}
	}

	return "", fmt.Errorf("IP address %s is does not belong to this user.", ipAddress)
}

func (i *InstanceUtils) CreateServer() error {
	name := i.driver.Name
	if name == "" {
		name = namegenerator.GetRandomName("docker-machine")
	}

	// IP address handling.
	ipAddress := i.driver.IPAddress
	dynamicIPRequired := i.driver.IPAddress != ""

	// FIXME: Remove this check and let the API handle the error?
	if !dynamicIPRequired && validation.IsUUID(ipAddress) {
		ID, err := i.getIPIDFromAddress(ipAddress)
		if err != nil {
			return err
		}
		ipAddress = ID
	}

	req := &instance.CreateServerRequest{
		Zone:              i.driver.Zone,
		Name:              name,
		DynamicIPRequired: scw.BoolPtr(dynamicIPRequired),
		CommercialType:    i.driver.Type,
		Image:             i.driver.Image,
		EnableIPv6:        i.driver.EnableIPV6,
		Organization:      i.driver.OrganizationID,
		Tags:              i.driver.Tags,
	}
	if ipAddress != "" {
		req.PublicIP = scw.StringPtr(ipAddress)
	}
	if i.driver.SecurityGroupID != "" {
		req.SecurityGroup = scw.StringPtr(i.driver.SecurityGroupID)
	}
	if i.driver.PlacementGroupID != "" {
		req.PlacementGroup = scw.StringPtr(i.driver.PlacementGroupID)
	}

	res, err := i.instanceAPI.CreateServer(req)
	if err != nil {
		return err
	}

	i.serverID = res.Server.ID

	return nil
}

func (i *InstanceUtils) GetCreatedServer() (*instance.Server, error) {
	res, err := i.instanceAPI.GetServer(&instance.GetServerRequest{ServerID: i.serverID})
	if err != nil {
		return nil, err
	}

	return res.Server, err
}

func (i *InstanceUtils) RemoveServer() error {
	_, err := i.instanceAPI.ServerAction(&instance.ServerActionRequest{
		ServerID: i.serverID,
		Action:   instance.ServerActionTerminate,
	})
	if err != nil {
		return fmt.Errorf("cannot remove the server: %s", err)
	}

	return nil
}

/*
- SCW_DEFAULT_ORGANIZATION_ID
- SCW_DEFAULT_ZONE
*/

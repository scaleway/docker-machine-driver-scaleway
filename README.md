<p align="center">
  <img src="https://raw.githubusercontent.com/scaleway/docker-machine-driver-scaleway/master/misc/logo_readme.png" width="800"/>
</p>

## Overview

[![Build
Status](https://travis-ci.org/scaleway/docker-machine-driver-scaleway.svg?branch=master)](https://travis-ci.org/scaleway/docker-machine-driver-scaleway)

A 3rd-party driver plugin for Docker machine to manage your containers on the servers of Scaleway

## Setup

```
# install docker-machine-driver-scaleway in your $GOPATH/bin
$> go get -u github.com/scaleway/docker-machine-driver-scaleway
```

## Usage

### 1. Get your Scaleway credentials

You can find your `ACCESS KEY` and generate your `TOKEN` [here](https://cloud.scaleway.com/#/credentials)

### 2. Scaleway driver helper
```console
$> docker-machine create -d scaleway -h
Usage: docker-machine create [OPTIONS] [arg...]

Create a machine

Description:
   Run 'docker-machine create --driver name' to include the create flags for that driver in the help text.

Options:

   --driver, -d "none"                                                                               Driver to create machine with.
   --engine-env [--engine-env option --engine-env option]                                            Specify environment variables to set in the engine
   --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]  Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"                                                     Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
   --engine-label [--engine-label option --engine-label option]                                      Specify labels for the created engine
   --engine-opt [--engine-opt option --engine-opt option]                                            Specify arbitrary flags to include with the created engine in the form flag=value
   --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]        Specify registry mirrors to use
   --engine-storage-driver                                                                           Specify a storage driver to use with the engine
   --scaleway-commercial-type "VC1S"                                                                 Specifies the commercial type [$SCALEWAY_COMMERCIAL_TYPE]
   --scaleway-name                                                                                   Assign a name [$SCALEWAY_NAME]
   --scaleway-organization                                                                           Scaleway organization [$SCALEWAY_ORGANIZATION]
   --scaleway-token                                                                                  Scaleway token [$SCALEWAY_TOKEN]
   --swarm                                                                                           Configure Machine with Swarm
   --swarm-addr                                                                                      addr to advertise for Swarm (default: detect and use the machine IP)
   --swarm-discovery                                                                                 Discovery service to use with Swarm
   --swarm-host "tcp://0.0.0.0:3376"                                                                 ip/socket to listen on for Swarm master
   --swarm-image "swarm:latest"                                                                      Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
   --swarm-master                                                                                    Configure Machine to be a Swarm master
   --swarm-opt [--swarm-opt option --swarm-opt option]                                               Define arbitrary flags for swarm
   --swarm-strategy "spread"                                                                         Define a default scheduling strategy for Swarm
   --tls-san [--tls-san option --tls-san option]                                                     Support extra SANs for TLS certs
```

### 3. Create your machine

Ensure you have your `ACCESS KEY` and a `TOKEN`

```
$> docker-machine create -d scaleway --scaleway-token=TOKEN --scaleway-organization=ACCESS_KEY --scaleway-name="cloud-scaleway-1" cloud-scaleway
Running pre-create checks...
Creating machine...
(cloud-scaleway) Creating SSH key...
(cloud-scaleway) Creating server...
(cloud-scaleway) Starting server...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(upstart)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env cloud-scaleway
```

### 4. Test your machine

```
$> eval $(docker-machine env cloud-scaleway)      # loads environment variables to use your machine

$> docker-machine ls                              # cloud-scaleway is now activated
NAME             ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
cloud-scaleway   *        scaleway     Running   tcp://212.47.248.251:2376           v1.10.3
dev              -        virtualbox   Running   tcp://192.168.99.100:2376           v1.9.1

$> docker run -d -p 80:80 owncloud:8.1            # starts a owncloud image
Unable to find image 'owncloud:8.1' locally
8.1: Pulling from library/owncloud
...

$> docker ps                                      # displays your containers
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                NAMES
ebdd86fcd18b        owncloud:8.1        "/entrypoint.sh apach"   22 seconds ago      Up 20 seconds       0.0.0.0:80->80/tcp   elegant_shirley

$> curl --silent http://212.47.248.251 | head -n1 # you can also open your browser with your IP
<!DOCTYPE html>
```

## Options

|Option Name                                                   |Description        |Default Value|required|
|--------------------------------------------------------------|-------------------|-------------|--------|
|``--scaleway-organization`` or ``$SCALEWAY_ORGANIZATION``     |Organization UUID  |none         |yes     |
|``--scaleway-token`` or ``$SCALEWAY_TOKEN``                   |Token UUID         |none         |yes     |
|``--scaleway-name`` or ``$SCALEWAY_NAME``                     |Server name        |none         |no      |
|``--scaleway-commercial-type`` or ``$SCALEWAY_COMMERCIAL_TYPE`` |Commercial type    |VC1S         |no      |

---

## Changelog

### v1.0.0 (2016-04-19)

* Sleep only when we stop an host ([#4](https://github.com/scaleway/docker-machine-driver-scaleway/issues/4))
* Loads credentials from `~/.scwrc` if available ([#2](https://github.com/scaleway/docker-machine-driver-scaleway/issues/2))
* Support of `create`
* Support of `start`
* Support of `stop`
* Support of `rm`
* Support of `restart`
* Support of `--scaleway-commercial-type`
* Support of `--scaleway-name`


---

## Development

Feel free to contribute :smiley::beers:

## Links

- **Scaleway console**: https://cloud.scaleway.com/
- **Scaleway cli**: https://github.com/scaleway/scaleway-cli
- **Scaleway github**: https://github.com/scaleway
- **Scaleway github-community**: https://github.com/scaleway-community
- **Docker Machine**: https://docs.docker.com/machine/
- **Report bugs**: https://github.com/scaleway/docker-machine-driver-scaleway/issues

## License

[MIT](https://github.com/scaleway/docker-machine-driver-scaleway/blob/master/LICENSE)

<p align="center">
  <img src="https://raw.githubusercontent.com/scaleway/docker-machine-driver-scaleway/master/misc/logo_readme.png" width="800"/>
</p>

## Overview

[![Build
Status](https://travis-ci.org/scaleway/docker-machine-driver-scaleway.svg?branch=master)](https://travis-ci.org/scaleway/docker-machine-driver-scaleway)

A 3rd-party driver plugin for Docker machine to manage your containers on the servers of Scaleway

## Setup

### Homebrew

```shell
# install latest release of docker-machine-driver-scaleway and docker-machine using homebrew
$ brew tap scaleway/scaleway
$ brew install scaleway/scaleway/docker-machine-driver-scaleway

# install latest (git) version of docker-machine-driver-scaleway
$ brew tap scaleway/scaleway
$ brew install scaleway/scaleway/docker-machine-driver-scaleway --HEAD
```

### Go
```shell
# install latest (git) version of docker-machine-driver-scaleway in your $GOPATH/bin (depends on Golang and docker-machine)
$ go get -u github.com/scaleway/docker-machine-driver-scaleway
```

### Binary

You can find sources and pre-compiled binaries [here](https://github.com/scaleway/docker-machine-driver-scaleway/releases/latest)

```shell
# Download the binary (this example downloads the binary for darwin amd64)
$ curl -sL https://github.com/scaleway/docker-machine-driver-scaleway/releases/download/v1.2.1/docker-machine-driver-scaleway_1.2.1_darwin_amd64.zip -O
$ unzip docker-machine-driver-scaleway_1.2.1_darwin_amd64.zip

# Make it executable and copy the binary in a directory accessible with your $PATH
$ chmod +x docker-machine-driver-scaleway_1.2.1_darwin_amd64/docker-machine-driver-scaleway
$ sudo cp docker-machine-driver-scaleway_1.2.1_darwin_amd64/docker-machine-driver-scaleway /usr/local/bin/
```

## Usage

### 1. Get your Scaleway credentials

You can find your `ACCESS KEY` and generate your `TOKEN` [here](https://cloud.scaleway.com/#/credentials)

### 2. Scaleway driver helper
```console
$ docker-machine create -d scaleway -h
Usage: docker-machine create [OPTIONS] [arg...]

Create a machine

Description:
   Run 'docker-machine create --driver name' to include the create flags for that driver in the help text.

Options:

   --driver, -d "none"                                                                               Driver to create machine with. [$MACHINE_DRIVER]
   --engine-env [--engine-env option --engine-env option]                                            Specify environment variables to set in the engine
   --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]  Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"                                                     Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
   --engine-label [--engine-label option --engine-label option]                                      Specify labels for the created engine
   --engine-opt [--engine-opt option --engine-opt option]                                            Specify arbitrary flags to include with the created engine in the form flag=value
   --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]        Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]
   --engine-storage-driver                                                                           Specify a storage driver to use with the engine
   --scaleway-commercial-type "VC1S"                                                                 Specifies the commercial type [$SCALEWAY_COMMERCIAL_TYPE]
   --scaleway-debug                                                                                  Enables Scaleway client debugging [$SCALEWAY_DEBUG]
   --scaleway-image "ubuntu-xenial"                                                                  Specifies the image [$SCALEWAY_IMAGE]
   --scaleway-bootscript "docker"                                                                    Specifies the bootscript [$SCALEWAY_BOOTSCRIPT]
   --scaleway-ip                                                                                     Specifies the IP address [$SCALEWAY_IP]
   --scaleway-ipv6                                                                                   Enable ipv6 [$SCALEWAY_IPV6]
   --scaleway-name                                                                                   Assign a name [$SCALEWAY_NAME]
   --scaleway-organization                                                                           Scaleway organization [$SCALEWAY_ORGANIZATION]
   --scaleway-port "22"                                                                              Specifies SSH port [$SCALEWAY_PORT]
   --scaleway-region "par1"                                                                          Specifies the location (par1,ams1) [$SCALEWAY_REGION]
   --scaleway-token                                                                                  Scaleway token [$SCALEWAY_TOKEN]
   --scaleway-user "root"                                                                            Specifies SSH user name [$SCALEWAY_USER]
   --scaleway-volumes                                                                                Attach additional volume (e.g., 50G) [$SCALEWAY_VOLUMES]
   --swarm                                                                                           Configure Machine to join a Swarm cluster
   --swarm-addr                                                                                      addr to advertise for Swarm (default: detect and use the machine IP)
   --swarm-discovery                                                                                 Discovery service to use with Swarm
   --swarm-experimental                                                                              Enable Swarm experimental features
   --swarm-host "tcp://0.0.0.0:3376"                                                                 ip/socket to listen on for Swarm master
   --swarm-image "swarm:latest"                                                                      Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
   --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]                                Define arbitrary flags for Swarm join
   --swarm-master                                                                                    Configure Machine to be a Swarm master
   --swarm-opt [--swarm-opt option --swarm-opt option]                                               Define arbitrary flags for Swarm master
   --swarm-strategy "spread"                                                                         Define a default scheduling strategy for Swarm
   --tls-san [--tls-san option --tls-san option]                                                     Support extra SANs for TLS certs
```

### 3. Create your machine

You need to configure your `ACCESS_KEY` and `TOKEN`, we suggest you to install [scw](https://github.com/scaleway/scaleway-cli) and create a credential file using `scw login`.

In the following example, authentication is done without any other dependencies using the `--scaleway-token=TOKEN` and `--scaleway-organization=ACCESS_KEY` parameters.

```console
$ docker-machine create -d scaleway --scaleway-token=TOKEN --scaleway-organization=ACCESS_KEY --scaleway-name="cloud-scaleway-1" cloud-scaleway
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

Note that you can store these parameters in the environment variables `SCALEWAY_TOKEN` and `SCALEWAY_ORGANIZATION`.

### 4. Test your machine

```console
$ eval $(docker-machine env cloud-scaleway)      # loads environment variables to use your machine

$ docker-machine ls                              # cloud-scaleway is now activated
NAME             ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
cloud-scaleway   *        scaleway     Running   tcp://212.47.248.251:2376           v1.10.3
dev              -        virtualbox   Running   tcp://192.168.99.100:2376           v1.9.1

$ docker run -d -p 80:80 owncloud:8.1            # starts a owncloud image
Unable to find image 'owncloud:8.1' locally
8.1: Pulling from library/owncloud
...

$ docker ps                                      # displays your containers
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                NAMES
ebdd86fcd18b        owncloud:8.1        "/entrypoint.sh apach"   22 seconds ago      Up 20 seconds       0.0.0.0:80->80/tcp   elegant_shirley

$ curl --silent http://212.47.248.251 | head -n1 # you can also open your browser with your IP
<!DOCTYPE html>
```

## Options

|Option Name                                                     |Description              |Default Value |required|
|----------------------------------------------------------------|-------------------------|--------------|--------|
|``--scaleway-organization`` or ``$SCALEWAY_ORGANIZATION``       |Organization UUID        |none          |yes     |
|``--scaleway-token`` or ``$SCALEWAY_TOKEN``                     |Token UUID               |none          |yes     |
|``--scaleway-name`` or ``$SCALEWAY_NAME``                       |Server name              |none          |no      |
|``--scaleway-commercial-type`` or ``$SCALEWAY_COMMERCIAL_TYPE`` |Commercial type          |VC1S          |no      |
|``--scaleway-image`` or ``$SCALEWAY_IMAGE``                     |Server image             |ubuntu-xenial |no      |
|``--scaleway-region`` or ``$SCALEWAY_REGION``                   |Specify the location     |par1          |no      |
|``--scaleway-debug`` or ``$SCALEWAY_DEBUG``                     |Toggle debugging         |false         |no      |
|``--scaleway-ip`` or ``$SCALEWAY_IP``                           |Server IP                |""            |no      |
|``--scaleway-volumes`` or ``$SCALEWAY_VOLUMES``                 |Attach additional volume |""            |no      |
|``--scaleway-user`` or ``$SCALEWAY_USER``                       |SSH User                 |root          |no      |
|``--scaleway-port`` or ``$SCALEWAY_PORT``                       |SSH port                 |22            |no      |

---

## Examples

```bash
# create a Scaleway docker host
docker-machine create -d scaleway my-scaleway-docker-machine

# create a VC1M server, name it my-docker-machine-1 on Scaleway and my-docker1 in the local Docker machine, with debug enabled
docker-machine create -d scaleway \
  --scaleway-name="my-docker-machine-1" --scaleway-debug \
  --scaleway-commercial-type="VC1M" --scaleway-volumes="50G" \
  my-docker1

# create a swarm master on a VC1M
docker-machine create -d scaleway \
  --scaleway-commercial-type="VC1M" --scaleway-volumes="50G" \
  --swarm --swarm-master --swarm-discovery="XXX"
  my-swarm-manager

# create a swarm slave on a VC1S
docker-machine create -d scaleway \
  --scaleway-commercial-type="VC1S" \
  --swarm --swarm-discovery="XXX"
  my-swarm-node

# create a docker host on the different server offers
docker-machine create -d scaleway --scaleway-commercial-type="VC1S"                           my-vc1s-node
docker-machine create -d scaleway --scaleway-commercial-type="VC1M" --scaleway-volumes="50G"  my-vc1m-node
docker-machine create -d scaleway --scaleway-commercial-type="VC1L" --scaleway-volumes="100G" my-vc1l-node
docker-machine create -d scaleway --scaleway-commercial-type="C2S"                            my-c2s-node
docker-machine create -d scaleway --scaleway-commercial-type="C2M"                            my-c2m-node
docker-machine create -d scaleway --scaleway-commercial-type="C2L"                            my-c2l-node
```

More [examples](https://github.com/scaleway/docker-machine-driver-scaleway/tree/master/examples).

---

## How to start an ARM server

To launch an ARM server, you need to start a server with our Docker Image, and use an empty install script.

```console
$ curl -sL http://bit.ly/1sf3j8V
#!/bin/sh

exit 0

$ docker-machine create -d scaleway --scaleway-commercial-type=C1 --scaleway-image=docker --engine-install-url="http://bit.ly/1sf3j8V" arm-machine
Running pre-create checks...
Creating machine...
(arm-machine) Creating SSH key...
(arm-machine) Creating server...
(arm-machine) Starting server...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env arm-machine

$ eval $(docker-machine env arm-machine) # arm-machine is now activated

$ docker run -it --rm multiarch/ubuntu-core:armhf-xenial # test an ARM container
Unable to find image 'multiarch/ubuntu-core:armhf-xenial' locally
armhf-xenial: Pulling from multiarch/ubuntu-core
9d12e3a67364: Pull complete
441bb0ba1886: Pull complete
4d9398209a87: Pull complete
89c0bb260a76: Pull complete
Digest: sha256:9b01beb4cdf0e1814583113105965f6b82a2fa618f403075f5ff653ac797911b
Status: Downloaded newer image for multiarch/ubuntu-core:armhf-xenial

root@ab197ef8bd3c:/# uname -a
Linux ab197ef8bd3c 4.5.4-docker-1 #1 SMP Thu May 19 18:02:43 UTC 2016 armv7l armv7l armv7l GNU/Linux
root@ab197ef8bd3c:/# exit
```


---

## Changelog

### v1.3 (2016-10-28)

* Add `--scaleway-region` to start server on different location e.g. `ams1` (Amsterdam)
* Fix `user-agent` format
* Add `--scaleway-ipv6` ([#50](https://github.com/scaleway/docker-machine-driver-scaleway/issues/50))
* Add `--scaleway-port`
* Add `--scaleway-user`

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.2.1...v1.3)

### v1.2.1 (2016-05-20)

* Delete IP only when she has been created by docker-machine

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.2.0...v1.2.1)

### v1.2.0 (2016-05-08)

* Add `--scaleway-volumes` ([#37](https://github.com/scaleway/docker-machine-driver-scaleway/issues/37))
* Add `--scaleway-ip` ([#40](https://github.com/scaleway/docker-machine-driver-scaleway/issues/40))

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.1.0...v1.2.0)

### v1.1.0 (2016-04-28)

* Fix provisionning error with xenial
* `docker-machine ls` displays the commercial-type
* Switch default image to **Ubuntu Xenial**
* Add `--scaleway-image` ([#22](https://github.com/scaleway/docker-machine-driver-scaleway/issues/22))
* Add `--scaleway-debug`

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.0.2...v1.1.0)

### v1.0.2 (2016-04-20)

* Add GOXC configuration ([#19](https://github.com/scaleway/docker-machine-driver-scaleway/issues/19))
* Fix rm subcommand ([#17](https://github.com/scaleway/docker-machine-driver-scaleway/issues/17))
* Initial homebrew support ([#9](https://github.com/scaleway/docker-machine-driver-scaleway/issues/9))

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.0.1...v1.0.2)

### v1.0.1 (2016-04-19)

* Bump dependencies

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.0.0...v1.0.1)

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

## Debugging

```console
$ SCALEWAY_DEBUG=1 MACHINE_DEBUG=1 docker-machine ...
```

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

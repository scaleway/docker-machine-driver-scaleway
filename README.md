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

<!-- ### Binary

You can find sources and pre-compiled binaries [here](https://github.com/scaleway/docker-machine-driver-scaleway/releases/latest)

```shell
# Download the binary (this example downloads the binary for darwin amd64)
$ curl -sL https://github.com/scaleway/docker-machine-driver-scaleway/releases/download/v1.2.1/docker-machine-driver-scaleway-v2_2.0.0_darwin_amd64.zip -O
$ unzip docker-machine-driver-scaleway-v2_2.0.0_darwin_amd64.zip

# Make it executable and copy the binary in a directory accessible with your $PATH
$ chmod +x docker-machine-driver-scaleway-v2_2.0.0_darwin_amd64/docker-machine-driver-scaleway
$ sudo cp docker-machine-driver-scaleway-v2_2.0.0_darwin_amd64/docker-machine-driver-scaleway /usr/local/bin/
``` -->

## Usage

### 1. Get your Scaleway credentials

You can find your `ACCESS KEY`, `SECRET KEY` and `PROJECT_ID` in section [API Keys](https://console.scaleway.com/project/credentials)

### 2. Scaleway driver helper
```console
$ docker-machine create -d scaleway-v2 --help
Usage: docker-machine create [OPTIONS] [arg...]

Create a machine

Description:
  Run 'docker-machine create --driver name --help' to include the create flags for that driver in the help text.

Options:

  --driver, -d "virtualbox"                                                                            Driver to create machine with. [$MACHINE_DRIVER]
  --engine-env [--engine-env option --engine-env option]                                               Specify environment variables to set in the engine
  --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]     Specify insecure registries to allow with the created engine
  --engine-install-url "https://get.docker.com"                                                        Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
  --engine-label [--engine-label option --engine-label option]                                         Specify labels for the created engine
  --engine-opt [--engine-opt option --engine-opt option]                                               Specify arbitrary flags to include with the created engine in the form flag=value
  --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]           Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]
  --engine-storage-driver                                                                              Specify a storage driver to use with the engine
  --scaleway-accesskey                                                                                 Scaleway accesskey (required) [$SCALEWAY_ACCESSKEY]
  --scaleway-bootscript                                                                                Specifies the bootscript [$SCALEWAY_BOOTSCRIPT]
  --scaleway-commercial-type "DEV1-S"                                                                  Specifies the commercial type [$SCALEWAY_COMMERCIAL_TYPE]
  --scaleway-debug                                                                                     Enables Scaleway client debugging [$SCALEWAY_DEBUG]
  --scaleway-image "ubuntu-focal"                                                                      Specifies the image [$SCALEWAY_IMAGE]
  --scaleway-ip                                                                                        Specifies the Public IP address [$SCALEWAY_IP]
  --scaleway-ipv6                                                                                      Enable ipv6 [$SCALEWAY_IPV6]
  --scaleway-name                                                                                      Assign a name [$SCALEWAY_NAME]
  --scaleway-project-id                                                                                Scaleway project id (required) [$SCALEWAY_PROJECT_ID]
  --scaleway-secretkey                                                                                 Scaleway secretkey (required) [$SCALEWAY_SECREYKEY]
  --scaleway-zone "fr-par-2"                                                                           Specifies the location (fr-par-1, fr-par-2, nl-ams-1, pl-waw-1) [$SCALEWAY_ZONE]
  --swarm                                                                                              Configure Machine to join a Swarm cluster
  --swarm-addr                                                                                         addr to advertise for Swarm (default: detect and use the machine IP)
  --swarm-discovery                                                                                    Discovery service to use with Swarm
  --swarm-experimental                                                                                 Enable Swarm experimental features
  --swarm-host "tcp://0.0.0.0:3376"                                                                    ip/socket to listen on for Swarm master
  --swarm-image "swarm:latest"                                                                         Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
  --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]                                   Define arbitrary flags for Swarm join
  --swarm-master                                                                                       Configure Machine to be a Swarm master
  --swarm-opt [--swarm-opt option --swarm-opt option]                                                  Define arbitrary flags for Swarm master
  --swarm-strategy "spread"                                                                            Define a default scheduling strategy for Swarm
  --tls-san [--tls-san option --tls-san option]                                                        Support extra SANs for TLS certs
```

### 3. Create your machine

```console
$ docker-machine create -d scaleway-v2 --scaleway-accesskey ACCESS_KEY --scaleway-secretkey SECRET_KEY --scaleway-project-id PROJECT_ID --scaleway-name cloud-scaleway-1 cloud-scaleway
Running pre-create checks...
Creating machine...
(cloud-scaleway) Create Scaleway Server ...
(cloud-scaleway) Creating SSH key...
(cloud-scaleway) Server created: cloud-scaleway-1 (<INSTANCE UUID>)
(cloud-scaleway) Start Server ...
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
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env cloud-scaleway
```

### 4. Test your machine

```console
$ eval $(docker-machine env cloud-scaleway)      # loads environment variables to use your machine

$ docker-machine ls
NAME             ACTIVE   DRIVER             STATE     URL                         SWARM   DOCKER      ERRORS
cloud-scaleway   -        scaleway(DEV1-S)   Running   tcp://IP_ADDRESS:2376           v19.03.13

$ docker run -d -p 80:80 owncloud:8.1            # starts a owncloud image
Unable to find image 'owncloud:8.1' locally
8.1: Pulling from library/owncloud
...

$ docker ps                                      # displays your containers
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                NAMES
ebdd86fcd18b        owncloud:8.1        "/entrypoint.sh apach"   22 seconds ago      Up 20 seconds       0.0.0.0:80->80/tcp   elegant_shirley

$ curl --silent http://IP_ADDRESS | head -n1 # you can also open your browser with your IP
<!DOCTYPE html>
```

## Options

|Option Name                                                     |Description              |Default Value |required|
|----------------------------------------------------------------|-------------------------|--------------|--------|
|``--scaleway-porject-id`` or ``$SCALEWAY_PROJECT_ID``           |Project UUID             |none          |yes     |
|``--scaleway-accesskey`` or ``$SCALEWAY_ACCESSKEY``             |Access Key UUID          |none          |yes     |
|``--scaleway-secretkey`` or ``$SCALEWAY_SECREYKEY``             |Secret Key UUID          |none          |yes     |
|``--scaleway-name`` or ``$SCALEWAY_NAME``                       |Server name              |none          |no      |
|``--scaleway-commercial-type`` or ``$SCALEWAY_COMMERCIAL_TYPE`` |Commercial type          |DEV1-S        |no      |
|``--scaleway-image`` or ``$SCALEWAY_IMAGE``                     |Server image             |ubuntu-focal  |no      |
|``--scaleway-zone`` or ``$SCALEWAY_ZONE``                       |Specify the location     |fr-par-2      |no      |
|``--scaleway-debug`` or ``$SCALEWAY_DEBUG``                     |Toggle debugging         |false         |no      |
|``--scaleway-ip`` or ``$SCALEWAY_IP``                           |Server IP                |""            |no      |
<!-- |``--scaleway-volumes`` or ``$SCALEWAY_VOLUMES``                 |Attach additional volume |""            |no      | -->

---

## Examples

```bash
# create a Scaleway docker host
docker-machine create -d scaleway my-scaleway-docker-machine

# create a DEV1-M server, name it my-docker-machine-1 on Scaleway and my-docker1 in the local Docker machine, with debug enabled
docker-machine --debug create -d scaleway \
  --scaleway-name="my-docker-machine-1" --scaleway-debug \
  --scaleway-commercial-type="DEV1-M" \
  my-docker1

# create a swarm master on a DEV1-M
docker-machine create -d scaleway \
  --scaleway-commercial-type="DEV1-M" \
  --swarm --swarm-master --swarm-discovery="XXX"
  my-swarm-manager

# create a swarm slave on a DEV1-M
docker-machine create -d scaleway \
  --scaleway-commercial-type="DEV1-M" \
  --swarm --swarm-discovery="XXX"
  my-swarm-node

# remove a machine
docker-machine rm cloud-scaleway
About to remove cloud-scaleway
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
(cloud-scaleway) Stop Server ...
(cloud-scaleway) Delete server: SERVER_UUID
(cloud-scaleway) Delete volume: VOLUME_UUID
(cloud-scaleway) Delete public IP: IP_ADDRESS (IP_UUID)
Successfully removed cloud-scaleway

# force remove a machine
docker-machine rm -f cloud-scaleway
About to remove cloud-scaleway
WARNING: This action will delete both local reference and remote instance.
(cloud-scaleway) Stop Server ...
(cloud-scaleway) Delete server: SERVER_UUID
(cloud-scaleway) Delete volume: VOLUME_UUID
(cloud-scaleway) Delete public IP: IP_ADDRESS (IP_UUID)
Successfully removed cloud-scaleway


```

More [examples](https://github.com/scaleway/docker-machine-driver-scaleway/tree/master/examples).

---

## Changelog

### v2.0.0 (2020-11-19)
* Use scaleway-sdk-go instead go-scaleway ( quick and dirty but functional :) ) 
* Fix remove instance with wrong state [issue #99](https://github.com/scaleway/docker-machine-driver-scaleway/issues/99)

### v1.6 (2018-12-03)

* Migrate from Godeps to dep
* Upgrade scaleway-cli dependency
* Add scripts to help release the project

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.5...v1.6)

### v1.5 (2018-11-19)

* Revert "Remove VC product line"
* Use xenial image id directly
* Use default image's bootscript

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.4...v1.5)

### v1.4 (2018-10-28)

* Change default bootscript
* Remove VC product line
* Allow the bootscript to be specified using it's unique id.
* Vendor update
* Remove IP adress if machine didn't exist ([#64](https://github.com/scaleway/docker-machine-driver-scaleway/pull/64))

View full [commits list](https://github.com/scaleway/docker-machine-driver-scaleway/compare/v1.3...v1.4)

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

* Add `--scaleway-ip` ([#37](https://github.com/scaleway/docker-machine-driver-scaleway/issues/37))
* Add `--scaleway-volumes` ([#40](https://github.com/scaleway/docker-machine-driver-scaleway/issues/40))

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

## Release

You can easily build for all supported platforms using `build-dmds-packages.sh` script
located in `./scripts`

## Links

- **Scaleway console**: https://cloud.scaleway.com/
- **Scaleway cli**: https://github.com/scaleway/scaleway-cli
- **Scaleway github**: https://github.com/scaleway
- **Scaleway github-community**: https://github.com/scaleway-community
- **Docker Machine**: https://docs.docker.com/machine/
- **Report bugs**: https://github.com/scaleway/docker-machine-driver-scaleway/issues

## License

[MIT](https://github.com/scaleway/docker-machine-driver-scaleway/blob/master/LICENSE)

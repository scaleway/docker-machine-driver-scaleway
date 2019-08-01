<p align="center"><img height="125" src="docs/static_files/scaleway-logo.png"/></p>

<p align="center">
  <a href="https://travis-ci.org/scaleway/docker-machine-driver-scaleway"><img src="https://travis-ci.org/scaleway/docker-machine-driver-scaleway.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://goreportcard.com/report/github.com/scalewaydocker-machine-driver-scaleway"><img src="https://goreportcard.com/badge/scaleway/docker-machine-driver-scaleway" alt="GoReportCard" /></a>
</p>

# Scaleway Docker Machine Driver

A 3rd-party driver plugin for Docker machine to manage your containers on Scaleway servers.

## Install

### Homebrew

Install the latest release using [homebrew](https://brew.sh/):

```bash
$ brew tap scaleway/scaleway
$ brew install scaleway/scaleway/docker-machine-driver-scaleway

# To install the HEAD version of this repository
$ brew install scaleway/scaleway/docker-machine-driver-scaleway --HEAD 
```

### Go

Install HEAD version in your `$GOPATH/bin` (depends on `Golang` and `docker-machine`)

```bash
$ go get -u github.com/scaleway/docker-machine-driver-scaleway
```

### Binary

You can find sources and pre-compiled binaries [here](https://github.com/scaleway/docker-machine-driver-scaleway/releases/latest).

Download the binary (this example downloads the binary for `darwin amd64`)

```bash
$ curl -sL https://github.com/scaleway/docker-machine-driver-scaleway/releases/download/v1.6/docker-machine-driver-scaleway_1.6_darwin_amd64.zip -O
$ unzip docker-machine-driver-scaleway_1.6_darwin_amd64.zip
```

Make the binary executable and copy it in a directory accessible with your `$PATH`:

```bash
$ chmod +x docker-machine-driver-scaleway_1.6_darwin_amd64/docker-machine-driver-scaleway
$ sudo cp docker-machine-driver-scaleway_1.6_darwin_amd64/docker-machine-driver-scaleway /usr/local/bin/
```

## Usage

At any time, you can read the driver helper with this command:

```bash
$ docker-machine create -d scaleway -h
```

### 1. Get your Scaleway credentials

The Scaleway authentication is based on an **organization ID** and a **secret key** (token).
You can find both of them in the section "API Tokens" of the [Scaleway Console](https://console.scaleway.com/account/credentials).

Since secret keys are only revealed one time (when it is first created) you might need to create a new one.
Click on the "Generate new token" button to create them. Giving it a friendly-name is recommended.

You can now set your environment variables:

```bash
export SCALEWAY_ORGANIZATION=<your-organization-id> # Node that you can also provide it your the --scaleway-organization flag
export SCALEWAY_TOKEN=<your-secret-key> # Node that you can also provide it your the --scaleway-token flag
```

### 2. Create your machine

```console
$ docker-machine create -d scaleway --scaleway-name="scw-machine" cloud-scaleway
Running pre-create checks...
Creating machine...
(cloud-scaleway) Creating SSH key...
(cloud-scaleway) Creating server...
(cloud-scaleway) Starting server...
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

$ docker-machine ls                              # cloud-scaleway is now activated
NAME             ACTIVE   DRIVER           STATE     URL                       SWARM   DOCKER     ERRORS
cloud-scaleway   *        scaleway(VC1S)   Running   tcp://51.158.119.9:2376           v19.03.1

$ docker run -d -p 80:80 owncloud:8.1            # starts a owncloud image
Unable to find image 'owncloud:8.1' locally
8.1: Pulling from library/owncloud
...

$ docker ps                                      # displays your containers
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                NAMES
a21475c67b10        owncloud:8.1        "/entrypoint.sh apac…"   8 seconds ago       Up 6 seconds        0.0.0.0:80->80/tcp   gallant_solomon

$ docker-machine ip cloud-scaleway               # get you machine public IP

$ curl --silent http://51.158.119.9 | head -n1   # you can also open your browser with your IP
<!DOCTYPE html>
```

## Options

|Flag or environment variable                                |Description               |Default Value   |required|
|------------------------------------------------------------|--------------------------|----------------|--------|
|`--scaleway-organization` or `$SCALEWAY_ORGANIZATION`       |Scaleway organization ID  |none            |yes     |
|`--scaleway-token` or `$SCALEWAY_TOKEN`                     |Scaleway secret key       |none            |yes     |
|`--scaleway-bootscript` or `$SCALEWAY_BOOTSCRIPT`           |Bootscript |none          |no              |no      |
|`--scaleway-commercial-type` or `$SCALEWAY_COMMERCIAL_TYPE` |Commercial type           |`VC1S`          |no      |
|`--scaleway-debug` or `$SCALEWAY_DEBUG`                     |Enables debug logs        |`false`         |no      |
|`--scaleway-image` or `$SCALEWAY_IMAGE`                     |Server image              |`ubuntu-xenial` |no      |
|`--scaleway-ip` or `$SCALEWAY_IP`                           |Server IP                 |""              |no      |
|`--scaleway-ipv6` or `$SCALEWAY_IP`                         |Enable IPv6               |""              |no      |
|`--scaleway-name` or `$SCALEWAY_NAME`                       |Server name               |none            |no      |
|`--scaleway-port` or `$SCALEWAY_PORT`                       |SSH port                  |`22`            |no      |
|`--scaleway-region` or `$SCALEWAY_REGION`                   |Specify the location      |`par1`          |no      |
|`--scaleway-user` or `$SCALEWAY_USER`                       |SSH User                  |`root`          |no      |
|`--scaleway-volumes` or `$SCALEWAY_VOLUMES`                 |Attach additional volumes |""              |no      |

---

## Examples

```console
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

# remove a machine
docker-machine rm my-vc1s-node
About to remove my-vc1s-node
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
Successfully removed my-vc1s-node

# force remove a machine
docker-machine rm -f my-vc1s-node
About to remove my-vc1s-node
WARNING: This action will delete both local reference and remote instance.
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

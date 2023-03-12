#!/bin/sh

docker-machine create -d scaleway --scaleway-commercial-type=C1   --scaleway-name=arm-small-docker arm-small
docker-machine create -d scaleway --scaleway-commercial-type=VC1S --scaleway-name=x64-small-docker x64-small
docker-machine create -d scaleway --scaleway-commercial-type=C2L  --scaleway-name=x64-large-docker x64-large

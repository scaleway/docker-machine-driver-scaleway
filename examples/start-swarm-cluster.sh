#!/bin/sh
# example of output: https://gist.github.com/moul/b8f604cac878a494dde123dc02d1863f


# generate a swarm unique token on Docker Hub discovery backend, for testing only
# doc: https://docs.docker.com/swarm/install-w-machine/#create-a-swarm-discovery-token
SWARM_TOKEN=$(docker run --rm swarm create)


( set -x
  # create swarm master
  docker-machine create \
                 -d scaleway \
                 --swarm \
                 --swarm-master \
                 --scaleway-name=swarm-manager \
                 --swarm-discovery=token://$SWARM_TOKEN swarm-manager
) &


# create 3 workers
for node in 1 2 3; do
    ( set -x
      docker-machine create \
                     -d scaleway \
                     --swarm \
                     --scaleway-name=swarm-node-$node \
                     --swarm-discovery=token://$SWARM_TOKEN \
                     swarm-node-$node
    ) &
done


# wait for background tasks
wait `jobs -p` || true


# configure shell for swarm
eval "$(docker-machine env --swarm swarm-manager)"


# list machines
docker-machine ls


# display infos
( set -x
  docker version
  docker info
)

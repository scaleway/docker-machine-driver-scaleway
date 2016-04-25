#!/bin/sh
## depends on docker client, curl, docker-machine + docker-machine-driver-scaleway, jq


# looking for existing rancher-server
if docker-machine inspect rancher-server >/dev/null; then
    echo "[-] rancher-server already exists."
    echo "    you need to remove your previous rancher machines from docker-machine before running this script."
    exit 1
fi


set -xe
SLAVES=3


# create a machine for the Rancher server
docker-machine create -d scaleway --scaleway-name=rancher-server rancher-server


# configure shell to use rancher-server
eval $(docker-machine env rancher-server)


# create machines for a Rancher host (slaves)
for i in $(seq 1 ${SLAVES}); do
    docker-machine create -d scaleway --scaleway-name=rancher-host-$i rancher-host-$i &
done


# start a Rancher Server on rancher-server
docker run -d --restart=always -p 8080:8080 rancher/server
(
    set +xe
    false
    while [ "$?" != "0" ]
    do
        # wait for Rancher Server installation
        docker logs $(docker ps -q) | grep "Startup Succeeded"
    done
)


# wait for background tasks
wait `jobs -p`; sleep 10



# generate a token
RANCHER_IP=$(docker-machine ip rancher-server)
curl "http://${RANCHER_IP}:8080/v1/registrationtoken" \
     -H 'x-api-project-id: 1a5' \
     -H 'Accept: application/json' \
     --data-binary '{"type":"registrationToken"}' \
     --compressed > /dev/null
sleep 5


# start a Rancher Agent on rancher-server
AGENT_CMD=$(
    curl "http://${RANCHER_IP}:8080/v1/registrationtokens?state=active&limit=-1" -H 'x-api-project-id: 1a5' -H 'Accept: application/json' --compressed \
        | jq -r '.data[0].command' \
        | sed 's/sudo//g'
         )
${AGENT_CMD}


# start a Rancher Agent on rancher-host-X (slaves)
for i in $(seq 1 ${SLAVES}); do
    eval $(docker-machine env rancher-host-$i)
    ${AGENT_CMD}
done


echo "Open your browser at http://${RANCHER_IP}:8080"

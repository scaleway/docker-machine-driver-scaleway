#!/bin/sh

# create a machine for the Rancher server
docker-machine create -d scaleway --scaleway-name=rancher-server rancher-server || exit 1

# configure shell to use rancher-server
eval $(docker-machine env rancher-server)

# create a machine for a Rancher host
docker-machine create -d scaleway --scaleway-name=rancher-host-1 rancher-host-1 || exit 1 &

# start a Rancher Server on rancher-server
docker run -d --restart=always -p 8080:8080 rancher/server || exit 1

false
while [ "$?" != "0" ]
do
	# wait for Rancher Server installation
	docker logs $(docker ps -q) | grep "Startup Succeeded"
done

# wait for background tasks
wait `jobs -p` || true

IP=$(docker-machine ip rancher-server)

sleep 10

# generate a token
curl 'http://'$IP':8080/v1/registrationtoken' -H 'x-api-project-id: 1a5' -H 'Accept: application/json' --data-binary '{"type":"registrationToken"}' --compressed > /dev/null
sleep 5
CMD=$(curl 'http://'$IP':8080/v1/registrationtokens?state=active&limit=-1' -H 'x-api-project-id: 1a5' -H 'Accept: application/json' --compressed | jq '.data[0].command' | sed 's/sudo//g' | tr -d '"')

# start a Rancher Agent on rancher-server
$CMD || exit 1

# configure shell to use rancher-host-1
eval $(docker-machine env rancher-host-1)

# start a Rancher Agent on rancher-host-1
$CMD || exit 1

echo "Open your browser at http://"$IP":8080"

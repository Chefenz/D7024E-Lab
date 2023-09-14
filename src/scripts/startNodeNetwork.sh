#!/bin/bash

#Traverse to the go main file
cd ./src/cmd/KademliaDDS/

#Build the go program 
echo "building the go program..."
env GOOS=linux GOARCH=amd64 go build main.go
echo "Done!"
echo ""

#Move the executable to the correct directory
mv main ../../bin

#Traverse back to the docker files
cd ./src

echo ""
echo "Starting up kademlia network and nodes"
echo ""

#Remove any running containers
echo "Remove any running containers"
sudo docker-compose down
echo ""

#Build the docker image and give it the correct name
echo "Build the docker file"
sudo docker build . -t kadlab
echo ""

#Start the docker containers
echo "start the docker containers"
sudo COMPOSE_HTTP_TIMEOUT=2500 docker-compose up -d
echo ""

#Check which containers are running
sudo docker ps
echo ""

echo "Done!"
#!/bin/bash

#Remove any running containers
echo "Remove any running containers"
sudo docker-compose down

#Build the docker image and give it the correct name
echo "Build the docker file"
sudo docker build . -t kadlab

#Start the docker containers
echo "start the docker containers"
sudo COMPOSE_HTTP_TIMEOUT=2500 docker-compose up -d

#Check which containers are running
sudo docker ps
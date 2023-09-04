#!/bin/bash

export COMPOSE_HTTP_TIMEOUT=200

#Build the docker image and give it the correct name
sudo docker build . -t kadlab

#Start the docker containers
sudo docker-compose up -d

#Check which containers are running
sudo docker-compose ps
#!/bin/bash

#Build the docker image and give it the correct name
sudo docker build . -t kadlab

#Start the docker containers
sudo COMPOSE_HTTP_TIMEOUT=2500 docker-compose up

#Check which containers are running
sudo docker-compose ps
#!/bin/bash

#Build the docker image and give it the correct name
sudo docker build . -t kadlab

echo "wait 30 seconds for the image to build"
sudo wait 30

#Start the docker containers
sudo docker-compose up -d

#Wait for the containers to start up
echo "wait 5 seconds for the containers to start up"
sudo wait 5

#Check which containers are running
sudo docker-compose ps
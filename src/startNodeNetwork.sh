#!/bin/bash

#Build the docker image and give it the correct name
sudo docker build . -t kadlab

#Start the docker containers
sudo docker-compose up -d

#Check which containers are running
sudo docker-compose ps
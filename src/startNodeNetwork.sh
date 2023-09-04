#!bin/bash

#Build the docker image and give it the correct name
docker build . -t kadlab

#Start the docker containers
docker-compose up -d

#Wait for the containers to start up
echo "wait ten seconds for the containers to start up"
wait 10

#Check which containers are running
docker-compose ps
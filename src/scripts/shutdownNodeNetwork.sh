#!/bin/bash

echo ""
echo "Shutting down kademlia network and nodes"
echo ""

#Shutdown and remove all running containers
sudo docker-compose down

echo "Done!"
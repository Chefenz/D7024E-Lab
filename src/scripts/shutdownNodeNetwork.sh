#!/bin/bash

echo ""
echo "Shutting down kademlia network and nodes"
echo ""

#Shutdown and remove all running containers
sudo COMPOSE_HTTP_TIMEOUT=2500 docker-compose down

echo "Done!"
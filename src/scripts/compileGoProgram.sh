#!/bin/bash

#Traverse to the go main file
cd ../cmd/KademliaDDS/

#Build the go program 
echo "building the go program..."
go build main.go
echo "Done!"
echo ""

#Move the executable to the correct directory
mv main.exe ../../bin
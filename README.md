# D7024E-Lab
Group 13 repository for lab solution

## Authors

Pontus Schünemann - ponsch-9@student.ltu.se

Emil Nyberg - eminyb-9@student.ltu.se

Alexander Österberg - aleste-9@student.ltu.se

## Dependencies
[Docker installation](https://docs.docker.com/engine/install/ubuntu/)

[Golang](https://go.dev/)

[Testify](https://github.com/stretchr/testify)

### Installation and setup of Docker
[Docker installation](https://docs.docker.com/engine/install/ubuntu/)


## How to run
To start the kademlia network use the script `startNodeNetwork.sh`. This will create one masternode and 49 regular nodes. The regular nodes will automatically add the master node and join its network. To run the startNodeNetwork script run `./startNodeNetwork.sh` in src/scripts.

If the program is changed and needs to be recompiled use the script `compileGoProgram.sh` by using the command `./compileGoProgram.sh` in src/scripts.

To shut down the node network use the script `shutdownNodeNetwork.sh` by using the command `./shutdownNodeNetwork.sh` in src/scripts

## Tests
The testing is done via [Testify](https://github.com/stretchr/testify) and to run the tests navigate to the folder you want to run the tests i.e **D7024E-Lab\src\labCode** and run the command `go test -cover`

## CLI
`get [hash]` - Returns a object corresponding to the inputted hash if it exists in the network

`put [string]` - Uploads a string to the network and returns a hash corresponding to the string

`exit` - Terminates the node

`togglePrints` - Toggles all non CLI interface printouts 

`help` - Provides information about a command

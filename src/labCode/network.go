package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Network struct {
	Me               Contact
	BucketChan       *chan Contact        // For updating buckets
	BucketWaitChan   *chan bool           // Wait for bucket update completion
	LookupChan       *chan Contact        // For lookup of contact
	FindChan         *chan Contact        // For find a contact
	ReturnFindChan   *chan []Contact      // For returning closest contacts to a contact
	DataReadChan     *chan ReadOperation  // For sending read requests to the data storage manager
	DataWriteChan    *chan WriteOperation // For sending write requests to the data storage manager
	CLIChan          *chan string         // For writing results to the command line interface
	FindConValueChan *chan FindContCloseToValOp
	rpcTimeOutChan   *chan bool
}

func NewNetwork(me Contact, bucketChan *chan Contact, bucketWaitChan *chan bool, lookupChan *chan Contact, findChan *chan Contact, returnFindChan *chan []Contact, dataReadChan *chan ReadOperation, dataWriteChan *chan WriteOperation, CLIChan *chan string, findContCloseToValOp *chan FindContCloseToValOp, rpcTimeOutChan *chan bool) Network {
	return Network{Me: me, BucketChan: bucketChan, BucketWaitChan: bucketWaitChan, LookupChan: lookupChan, FindChan: findChan, ReturnFindChan: returnFindChan, DataReadChan: dataReadChan, DataWriteChan: dataWriteChan, CLIChan: CLIChan, FindConValueChan: findContCloseToValOp, rpcTimeOutChan: rpcTimeOutChan}
}

func (network *Network) Listen(ip string, port int, stopChan chan string) {

	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	chk(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	chk(err)

	fmt.Println("Listening to: ", udpAddr)

	defer conn.Close()

	buffer := make([]byte, 4096)

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping listen..")
			conn.Close()
			return
		default:
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from UDP connection:", err)
				continue
			}
			if len(buffer) > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])

				go network.handleRPC(data)
			}
		}
	}
}

func (network *Network) handleRPC(data []byte) {

	// Unmarshal transfer data
	var transmitObj TransmitObj
	err := json.Unmarshal(data, &transmitObj)
	chk(err)

	//fmt.Println("Handling RPC: ", transmitObj.Message)

	switch transmitObj.Message {
	case "PING":
		contact := decodeTransmitObj(transmitObj, "Contact").(*Contact)
		*network.BucketChan <- *contact
		<-*network.BucketWaitChan

		transmitObj := TransmitObj{Message: "PONG", Data: network.Me}
		network.sendMessage(&transmitObj, contact)
	case "PONG":
		contact := decodeTransmitObj(transmitObj, "Contact").(*Contact)
		*network.BucketChan <- *contact
		<-*network.BucketWaitChan
	case "FIND_CONTACT":
		findContactPayload := decodeTransmitObj(transmitObj, "FindContactPayload").(*FindContactPayload)
		*network.BucketChan <- findContactPayload.Sender
		<-*network.BucketWaitChan

		// Lookup closest contacts over channels
		*network.FindChan <- findContactPayload.Target
		closestContacts := <-*network.ReturnFindChan

		returnFindContactPayload := ReturnFindContactPayload{Shortlist: closestContacts, Target: findContactPayload.Target}
		transmitObj := TransmitObj{Message: "RETURN_FIND_CONTACT", Data: returnFindContactPayload}

		network.sendMessage(&transmitObj, &findContactPayload.Sender)

	case "RETURN_FIND_CONTACT":
		returnFindContactPayload := decodeTransmitObj(transmitObj, "ReturnFindContactPayload").(*ReturnFindContactPayload)

		foundTarget := network.checkForFindContact(*returnFindContactPayload)

		if foundTarget == false {
			fmt.Println("Did Not Find The Target Node Will Try Again")
			*network.LookupChan <- returnFindContactPayload.Target
		}

	case "FIND_VALUE":
		findValuePayload := decodeTransmitObj(transmitObj, "FindValuePayload").(*FindValuePayload)
		sentFrom := transmitObj.Sender

		key := findValuePayload.Key

		requestRead := ReadOperation{Key: key.String(), Resp: make(chan []byte)}
		*network.DataReadChan <- requestRead

		result := <-requestRead.Resp

		//The data was found
		if result != nil {
			returnFindValueDataPayload := ReturnFindValuePayload{Data: string(result), Shortlist: nil, TargetKey: nil}
			transmitObj := TransmitObj{Message: "RETURN_FIND_VALUE", Sender: network.Me, Data: returnFindValueDataPayload}
			network.sendMessage(&transmitObj, &sentFrom)

		} else { // The data could not be found
			requestFindClosesTContactOp := FindContCloseToValOp{TargetID: key, Resp: make(chan []Contact)}

			*network.FindConValueChan <- requestFindClosesTContactOp
			closestContacts := <-requestFindClosesTContactOp.Resp

			returnFindValuePayload := ReturnFindValuePayload{Data: "", Shortlist: closestContacts, TargetKey: key}
			transmitObj := TransmitObj{Message: "RETURN_FIND_VALUE", Sender: network.Me, Data: returnFindValuePayload}
			network.sendMessage(&transmitObj, &sentFrom)

		}

	case "RETURN_FIND_VALUE":
		returnFindValueDataPayload := decodeTransmitObj(transmitObj, "ReturnFindValuePayload").(*ReturnFindValuePayload)

		targetID := returnFindValueDataPayload.TargetKey
		valueResult := returnFindValueDataPayload.Data
		closestContacts := returnFindValueDataPayload.Shortlist

		if valueResult != "" {
			// Check if dataResult is not empty
			select {
			case *network.CLIChan <- valueResult + " " + transmitObj.Sender.String():
				fmt.Println("I had the right result so I wrote", transmitObj.Sender.String())
			default:
				fmt.Println("I had the right result but someone already wrote so I skipped", transmitObj.Sender.String())
			}
		} else {
			// Check if dataResult is empty
			select {
			case <-*network.rpcTimeOutChan:
				fmt.Println("The RPC call timed out")
				*network.CLIChan <- "Could not find the result"
			default:
				fmt.Println("I did not have the valid result and it had not been posted so I send out another find value", transmitObj.Sender.String())

				findValuePayload := FindValuePayload{Key: targetID}
				transmitObj := TransmitObj{Message: "FIND_VALUE", Sender: network.Me, Data: findValuePayload}
				for i := 0; i < len(closestContacts); i++ {
					network.sendMessage(&transmitObj, &closestContacts[i])
				}
			}
		}

	case "STORE":
		storePayload := decodeTransmitObj(transmitObj, "StorePayload").(*StorePayload)
		sentFrom := transmitObj.Sender

		key := storePayload.Key
		dataStr := storePayload.Data
		data := []byte(dataStr)

		requestWrite := WriteOperation{Key: key.String(), Data: data, Resp: make(chan bool)}
		*network.DataWriteChan <- requestWrite

		<-requestWrite.Resp

		returnStorePayload := ReturnStorePayload{Key: key}
		transmitObj := TransmitObj{Message: "RETURN_STORE", Sender: network.Me, Data: returnStorePayload}
		network.sendMessage(&transmitObj, &sentFrom)

	case "RETURN_STORE":
		returnStorePayload := decodeTransmitObj(transmitObj, "ReturnStorePayload").(*ReturnStorePayload)

		key := returnStorePayload.Key

		select {
		case *network.CLIChan <- key.String():
			//fmt.Println("I WROTE")
		default:
			//fmt.Println("Someone already wrote the answer so I skipped")
		}
	}
}

func decodeTransmitObj(obj TransmitObj, objType string) interface{} {
	objMap, ok := obj.Data.(map[string]interface{})

	if ok != true {
		fmt.Println("Data is not a Map")
	}

	switch objType {
	case "Contact":
		var contact *Contact
		err := mapstructure.Decode(objMap, &contact)
		chk(err)
		return contact

	case "FindContactPayload":
		var findContactPayload *FindContactPayload
		err := mapstructure.Decode(objMap, &findContactPayload)
		chk(err)
		return findContactPayload

	case "ReturnFindContactPayload":
		var returnFindContactPayload *ReturnFindContactPayload
		err := mapstructure.Decode(objMap, &returnFindContactPayload)
		chk(err)
		return returnFindContactPayload

	case "StorePayload":
		var storePayload *StorePayload
		err := mapstructure.Decode(objMap, &storePayload)
		chk(err)
		return storePayload

	case "ReturnStorePayload":
		var returnStorePayload *ReturnStorePayload
		err := mapstructure.Decode(objMap, &returnStorePayload)
		chk(err)
		return returnStorePayload

	case "FindValuePayload":
		var findValuePayload *FindValuePayload
		err := mapstructure.Decode(objMap, &findValuePayload)
		chk(err)
		return findValuePayload

	case "ReturnFindValuePayload":
		var returnFindValuePayload *ReturnFindValuePayload
		err := mapstructure.Decode(objMap, &returnFindValuePayload)
		chk(err)
		return returnFindValuePayload
	}

	return nil

}

func (network *Network) checkForFindContact(returnFindContactPayload ReturnFindContactPayload) (foundTarget bool) {

	shortlist := returnFindContactPayload.Shortlist
	target := returnFindContactPayload.Target
	foundTarget = false

	for i := 0; i < len(shortlist); i++ {
		if *shortlist[i].ID != *network.Me.ID {
			contact := shortlist[i]
			*network.BucketChan <- contact
			<-*network.BucketWaitChan
		}

		if *shortlist[i].ID == *target.ID {
			foundTarget = true
			fmt.Println("Found The Target Node :)")
		}
	}

	return foundTarget

}

func (network *Network) sendMessage(transmitObj *TransmitObj, contact *Contact) {

	targetAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	chk(err)

	conn, err := net.DialUDP("udp", nil, targetAddr)
	chk(err)

	// Marshal the struct into JSON
	sendJSON, err := json.Marshal(transmitObj)
	chk(err)

	_, err = conn.Write(sendJSON)
	chk(err)

	conn.Close()

}

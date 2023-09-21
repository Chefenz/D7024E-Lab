package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Network struct {
	Me             Contact
	BucketChan     *chan Contact        // For update bucket
	BucketWaitChan *chan bool           // Wait for bucket update completion
	LookupChan     *chan Contact        // For lookup of contact
	FindChan       *chan Contact        // For find a contact
	ReturnFindChan *chan []Contact      // For returning closest contacts to a contact
	DataReadChan   *chan ReadOperation  //For sending read requests to the data storage
	DataWriteChan  *chan WriteOperation //For sending write requests to the data storage
}

func NewNetwork(me Contact, bucketChan *chan Contact, bucketWaitChan *chan bool, lookupChan *chan Contact, findChan *chan Contact, returnFindChan *chan []Contact, dataReadChan *chan ReadOperation, dataWriteChan *chan WriteOperation) Network {
	return Network{Me: me, BucketChan: bucketChan, BucketWaitChan: bucketWaitChan, LookupChan: lookupChan, FindChan: findChan, ReturnFindChan: returnFindChan, DataReadChan: dataReadChan, DataWriteChan: dataWriteChan}
}

func (network *Network) Listen(ip string, port int) {
	for newPort := port; newPort <= port+50; newPort++ {

		udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(newPort))
		chk(err)
		conn, err := net.ListenUDP("udp", udpAddr)
		chk(err)

		fmt.Println("Listening to: ", udpAddr)

		defer conn.Close()

		buffer := make([]byte, 4096)

		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from UDP connection:", err)
				continue
			}
			if len(buffer) > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])

				go network.handleRPC(data, conn)
			}
		}
	}
}

func (network *Network) handleRPC(data []byte, conn *net.UDPConn) {

	// Unmarshal transfer data
	var transmitObj TransmitObj
	err := json.Unmarshal(data, &transmitObj)
	chk(err)

	fmt.Println("Handling RPC: ", transmitObj.Message)

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
	case "HEARTBEAT":
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

	case "FIND_DATA":
		fmt.Println("This should handle finddata")
	case "STORE":
		storePayload := decodeTransmitObj(transmitObj, "StorePayload").(*StorePayload)
		sentFrom := transmitObj.Sender

		key := storePayload.Key
		dataStr := storePayload.Data
		data := []byte(dataStr)

		requestWrite := WriteOperation{Key: key.String(), Data: data, Resp: make(chan bool)}
		*network.DataWriteChan <- requestWrite

		returnStorePayload := ReturnStorePayload{Key: key}
		transmitObj := TransmitObj{Message: "RETURN_STORE", Sender: network.Me, Data: returnStorePayload}
		network.sendMessage(&transmitObj, &sentFrom)

	case "RETURN_STORE":
		returnStorePayload := decodeTransmitObj(transmitObj, "ReturnStorePayload").(*ReturnStorePayload)

		//Tell the waitgroup that it is done and wait for all the other store ndn retutn store to come trough
		fmt.Println("In return store after wait", returnStorePayload)

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

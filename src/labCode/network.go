package labCode

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Network struct {
	Me               Contact
	BucketChan       *chan Contact        // For update bucket
	BucketWaitChan   *chan bool           // Wait for bucket update completion
	LookupChan       *chan LookupContOp   // For lookup of contact
	FindChan         *chan Contact        // For find a contact
	ReturnFindChan   *chan []Contact      // For returning closest contacts to a contact
	DataReadChan     *chan ReadOperation  //For sending read requests to the data storage
	DataWriteChan    *chan WriteOperation //For sending write requests to the data storage
	CLIChan          *chan string
	FindConValueChan *chan FindContCloseToValOp
	FoundTarget      bool
	FoundValue       bool
}

func NewNetwork(me Contact, bucketChan *chan Contact, bucketWaitChan *chan bool, lookupChan *chan LookupContOp, findChan *chan Contact, returnFindChan *chan []Contact, dataReadChan *chan ReadOperation, dataWriteChan *chan WriteOperation, CLIChan *chan string, findContCloseToValOp *chan FindContCloseToValOp) Network {
	return Network{Me: me, BucketChan: bucketChan, BucketWaitChan: bucketWaitChan, LookupChan: lookupChan, FindChan: findChan, ReturnFindChan: returnFindChan, DataReadChan: dataReadChan, DataWriteChan: dataWriteChan, CLIChan: CLIChan, FindConValueChan: findContCloseToValOp, FoundTarget: false, FoundValue: false}
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

	fmt.Println("Handling RPC: ", transmitObj.Message)

	if time.Since(transmitObj.RPC_created_at) < 4*time.Second {
		fmt.Println("RPC not timed out with duration: ", time.Since(transmitObj.RPC_created_at))

		switch transmitObj.Message {
		case "PING":
			contact := decodeTransmitObj(transmitObj, "Contact").(*Contact)
			*network.BucketChan <- *contact
			<-*network.BucketWaitChan

			transmitObj := TransmitObj{Message: "PONG", Data: network.Me, RPC_created_at: transmitObj.RPC_created_at}
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
			transmitObj := TransmitObj{Message: "RETURN_FIND_CONTACT", Data: returnFindContactPayload, RPC_created_at: transmitObj.RPC_created_at}

			network.sendMessage(&transmitObj, &findContactPayload.Sender)

		case "RETURN_FIND_CONTACT":
			returnFindContactPayload := decodeTransmitObj(transmitObj, "ReturnFindContactPayload").(*ReturnFindContactPayload)

			network.checkForFindContact(*returnFindContactPayload)

			if network.FoundTarget == false {
				fmt.Println("Did Not Find The Target Node Will Try Again")
				lookupContOp := LookupContOp{Contact: &returnFindContactPayload.Target, RPC_created_at: transmitObj.RPC_created_at}
				*network.LookupChan <- lookupContOp
			}

		case "FIND_VALUE":
			findValuePayload := decodeTransmitObj(transmitObj, "FindValuePayload").(*FindValuePayload)
			sentFrom := transmitObj.Sender

			key := findValuePayload.Key

			requestRead := ReadOperation{Key: key.String(), Resp: make(chan []byte)}
			*network.DataReadChan <- requestRead

			result := <-requestRead.Resp
			//fmt.Println("The Result of the read operation in network:", result)

			if result != nil {
				returnFindValueDataPayload := ReturnFindValuePayload{Data: string(result), Shortlist: nil, TargetKey: nil}
				transmitObj := TransmitObj{Message: "RETURN_FIND_VALUE", Sender: network.Me, Data: returnFindValueDataPayload, RPC_created_at: transmitObj.RPC_created_at}
				network.sendMessage(&transmitObj, &sentFrom)

			} else {
				requestFindClosesTContactOp := FindContCloseToValOp{TargetID: key, Resp: make(chan []Contact)}

				*network.FindConValueChan <- requestFindClosesTContactOp
				closestContacts := <-requestFindClosesTContactOp.Resp

				returnFindValuePayload := ReturnFindValuePayload{Data: "", Shortlist: closestContacts, TargetKey: key}
				transmitObj := TransmitObj{Message: "RETURN_FIND_VALUE", Sender: network.Me, Data: returnFindValuePayload, RPC_created_at: transmitObj.RPC_created_at}
				network.sendMessage(&transmitObj, &sentFrom)

			}

		case "RETURN_FIND_VALUE":
			returnFindValueDataPayload := decodeTransmitObj(transmitObj, "ReturnFindValuePayload").(*ReturnFindValuePayload)

			network.checkForFindValue(transmitObj, *returnFindValueDataPayload)

			if network.FoundValue == false {
				fmt.Println("Did Not Find The Target Value Will Try Again")
				findValuePayload := FindValuePayload{Key: returnFindValueDataPayload.TargetKey}
				transmitObj := TransmitObj{Message: "FIND_VALUE", Sender: network.Me, Data: findValuePayload, RPC_created_at: transmitObj.RPC_created_at}
				for i := 0; i < len(returnFindValueDataPayload.Shortlist); i++ {
					network.sendMessage(&transmitObj, &returnFindValueDataPayload.Shortlist[i])
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
			transmitObj := TransmitObj{Message: "RETURN_STORE", Sender: network.Me, Data: returnStorePayload, RPC_created_at: transmitObj.RPC_created_at}
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
	} else {
		fmt.Println("RPC timed out with duration: ", time.Since(transmitObj.RPC_created_at))
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

func (network *Network) checkForFindContact(returnFindContactPayload ReturnFindContactPayload) {

	shortlist := returnFindContactPayload.Shortlist
	target := returnFindContactPayload.Target

	for i := 0; i < len(shortlist); i++ {
		if *shortlist[i].ID != *network.Me.ID {
			contact := shortlist[i]
			*network.BucketChan <- contact
			<-*network.BucketWaitChan
		}

		if *shortlist[i].ID == *target.ID {
			network.FoundTarget = true
			fmt.Println("Found The Target Node :)")
			timer := time.NewTimer(10 * time.Second)
			go func() {
				<-timer.C
				network.FoundTarget = false
			}()
		}
	}
}

func (network *Network) checkForFindValue(transmitObj TransmitObj, returnFindValueDataPayload ReturnFindValuePayload) {

	valueResult := returnFindValueDataPayload.Data

	if valueResult != "" {
		network.FoundValue = true
		*network.CLIChan <- valueResult + " " + transmitObj.Sender.String()
		timer := time.NewTimer(10 * time.Second)
		go func() {
			<-timer.C
			network.FoundValue = false
		}()
	}
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

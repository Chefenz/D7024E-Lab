package labCode

/*
import (
	"container/list"
)
*/

type Kademlia struct {
	routingTable RoutingTable
	network      Network
	//data         ToBeDetermined
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	count := 3
	shortList := kademlia.routingTable.FindClosestContacts(target.ID, count)

	for i := 0; i <= count-1; i++ {
		kademlia.network.SendFindContactMessage(&shortList[i])
		//kademlia.routingTable.AddContact(shortList[0]) lägg till contact från svar av SendFindContactMessage
	}

	if &shortList[0] == target { //Kan ändra shortlist till svar från findContact message
		kademlia.routingTable.AddContact(shortList[0])
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

/*
func JoinNetwork() Network {
	node := InitKademliaNode()
	network := Network{masterID}
	node.routingTable.AddContact(NewContact(&masterID, masterIP))
	//node.LookupContact()
	return network
}
*/

func InitKademliaNode() Kademlia {
	id := NewRandomKademliaID()
	ip := ""
	rt := NewRoutingTable(NewContact(id, ip))
	Listen(rt.me.Address, 8050)
	return Kademlia{*rt, nil}
}

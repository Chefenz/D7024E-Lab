package labCode

import "fmt"

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	Me             Contact
	buckets        [IDLength * 8]*bucket
	BucketChan     *chan Contact
	FindChan       *chan Contact
	ReturnFindChan *chan []Contact
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact, bucketChan *chan Contact, findChan *chan Contact, returnFindChan *chan []Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.Me = me
	routingTable.BucketChan = bucketChan
	routingTable.FindChan = findChan
	routingTable.ReturnFindChan = returnFindChan
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.Me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

func (routingTable *RoutingTable) UpdateBucketRoutine() {
	contact := <-*routingTable.BucketChan

	routingTable.AddContact(contact)

	fmt.Println("node has been updated in bucket")

}

func (routingTable *RoutingTable) FindClosestContactsRoutine() {
	target := <-*routingTable.FindChan

	closestContacts := routingTable.FindClosestContacts(target.ID, alpha)

	*routingTable.ReturnFindChan <- closestContacts
}

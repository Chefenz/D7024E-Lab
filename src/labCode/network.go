package labCode

type Network struct {
}

func NewNetwork() Network {
	return Network{}
}

func (network *Network) Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

type listenMessage struct {
	Message string
	Data    interface{}
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

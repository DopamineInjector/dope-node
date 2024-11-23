package messages

type NewConnectionMessage struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type MessageType struct {
	Type string `json:"type"`
}

type AvailableNodesAddresses struct {
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

type Transaction struct {
	Type     string  `json:"type"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

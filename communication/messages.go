package communication

type NewConnectionMessage struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type AvailableNodesAddresses struct {
	Addresses []string `json:"addresses"`
}

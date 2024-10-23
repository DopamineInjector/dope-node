package communication

type Message struct {
	Content string `json:"content"`
}

type FetchBootstrapAddressesResponse struct {
	Addresses []string
}

type ConnectRequest struct {
	sender string
}

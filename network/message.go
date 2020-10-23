package network

type MessageType int

type Message struct {
	endpoint string
	value    []byte
}

package common

type FailureCallback func(message *Message, err error)

type MessageHeader struct {
	Topic     string
	Offset    int64
	Key       string
	Partition int32
}

type Message struct {
	MessageHeader
	MessageAdditional
	Value []byte
}

type MessageAdditional struct {
	OnFailure FailureCallback
}

package message

import (
	"encoding/json"
)

type MozGenericMessageInterface interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) (MozGenericMessageInterface, error)
}

type publicMozGenericMessage struct {
	MozType 				string			`json:"type"`
	Id							int					`json:"id"`
	ThreadId				int					`json:"threadId"`
	// Body						string			`json:"body"`
	// Subject					string			`json:"subject"`
	// Smil						interface{}	`json:"smil"`
	// Attachments			interface{}	`json:"attachments"`
	// ExpiryDate			int					`json:"expiryDate"`
	Delivery				string			`json:"delivery"`
	DeliveryStatus	string			`json:"deliveryStatus"`
	Read						bool				`json:"read"`
	Receiver				string			`json:"receiver"`
	Sender					string			`json:"sender"`
	Timestamp				int					`json:"timestamp"`
	// MessageClass		string			`json:"class"`
}

type MozGenericMessage struct {
	mozType					string // sms or mms
	id							int
	threadId				int
	delivery				string
	deliveryStatus	string
	read						bool
	receiver				string
	sender					string
	timestamp				int
}

func (m *MozGenericMessage) GetType() string {
	return m.mozType
}

func (m *MozGenericMessage) setType(mozType string) string {
	m.mozType = mozType
	return m.mozType
}

func (m *MozGenericMessage) GetId() int {
	return m.id
}

func (m *MozGenericMessage) setId(id int) int {
	m.id = id
	return m.id
}

func (m *MozGenericMessage) GetThreadId() int {
	return m.threadId
}

func (m *MozGenericMessage) setThreadId(threadId int) int {
	m.threadId = threadId
	return m.threadId
}

func (m *MozGenericMessage) GetDelivery() string {
	return m.delivery
}

func (m *MozGenericMessage) setDelivery(delivery string) string {
	m.delivery = delivery
	return m.delivery
}

func (m *MozGenericMessage) GetDeliveryStatus() string {
	return m.deliveryStatus
}

func (m *MozGenericMessage) setDeliveryStatus(deliveryStatus string) string {
	m.deliveryStatus = deliveryStatus
	return m.deliveryStatus
}

func (m *MozGenericMessage) GetRead() bool {
	return m.read
}

func (m *MozGenericMessage) setRead(read bool) bool {
	m.read = read
	return m.read
}

func (m *MozGenericMessage) GetReceiver() string {
	return m.receiver
}

func (m *MozGenericMessage) setReceiver(receiver string) string {
	m.receiver = receiver
	return m.receiver
}

func (m *MozGenericMessage) GetSender() string {
	return m.sender
}

func (m *MozGenericMessage) setSender(sender string) string {
	m.sender = sender
	return m.sender
}

func (m *MozGenericMessage) GetTimestamp() int {
	return m.timestamp
}

func (m *MozGenericMessage) setTimestamp(timestamp int) int {
	m.timestamp = timestamp
	return m.timestamp
}

func (m MozGenericMessage) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (m MozGenericMessage) UnmarshalJSON(data []byte) (MozGenericMessageInterface, error) {
	var generic *publicMozGenericMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		return m, err
	}
	if (generic.MozType == "sms") {
		sms := &MozSmsMessage{}
		return sms.UnmarshalJSON(data)
	} else if (generic.MozType == "mms") {
		mms := &MozMmsMessage{}
		return mms.UnmarshalJSON(data)
	}
	return m, nil
}

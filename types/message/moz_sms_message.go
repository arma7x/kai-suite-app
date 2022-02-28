package message

import (
	"encoding/json"
)

type Class int

const (
	NORMAL	Class = iota
	CLASS_0
	CLASS_1
	CLASS_2
	CLASS_3
)

type publicMozSmsMessage struct {
	MozType 				string		`json:"type"`
	Id							int				`json:"id"`
	ThreadId				int				`json:"threadId"`
	Body						string		`json:"body"`
	Delivery				string		`json:"delivery"`
	DeliveryStatus	string		`json:"deliveryStatus"`
	Read						bool			`json:"read"`
	Receiver				string		`json:"receiver"`
	Sender					string		`json:"sender"`
	Timestamp				int				`json:"timestamp"`
	MessageClass		string		`json:"class"`
}

type MozSmsMessage struct {
	MozGenericMessage
	body						string
	messageClass		string
}

func (m *MozSmsMessage) GetBody() string {
	return m.body
}

func (m *MozSmsMessage) setBody(body string) string {
	m.body = body
	return m.body
}

func (m *MozSmsMessage) GetMessageClass() string {
	return m.messageClass
}

func (m *MozSmsMessage) setMessageClass(messageClass string) string {
	m.messageClass = messageClass
	return m.messageClass
}

func (m MozSmsMessage) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicMozSmsMessage{
		MozType: m.mozType,
		Id: m.id,
		ThreadId: m.threadId,
		Body: m.body,
		Delivery: m.delivery,
		DeliveryStatus: m.deliveryStatus,
		Read: m.read,
		Receiver: m.receiver,
		Sender: m.sender,
		Timestamp: m.timestamp,
		MessageClass: m.messageClass,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (m MozSmsMessage) UnmarshalJSON(data []byte) (MozGenericMessageInterface, error) {
	var public *publicMozSmsMessage
	if err := json.Unmarshal(data, &public); err != nil {
		return m, err
	}
	m.setType(public.MozType)
	m.setId(public.Id)
	m.setThreadId(public.ThreadId)
	m.setBody(public.Body)
	m.setDelivery(public.Delivery)
	m.setDeliveryStatus(public.DeliveryStatus)
	m.setRead(public.Read)
	m.setReceiver(public.Receiver)
	m.setSender(public.Sender)
	m.setTimestamp(public.Timestamp)
	m.setMessageClass(public.MessageClass)
	return m, nil
}

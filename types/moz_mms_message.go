package types

//import (
	//"encoding/json"
//)

//type publicMozMmsMessage struct {
	//MozType 				string			`json:"type"`
	//Id							int					`json:"id"`
	//ThreadId				int					`json:"threadId"`
	//Subject					string			`json:"subject"`
	//Smil						interface{}	`json:"smil"`
	//Attachments			interface{}	`json:"attachments"`
	//ExpiryDate			int					`json:"expiryDate"`
	//Delivery				string			`json:"delivery"`
	//DeliveryStatus	string			`json:"deliveryStatus"`
	//Read						bool				`json:"read"`
	//Receiver				string			`json:"receiver"`
	//Sender					string			`json:"sender"`
	//Timestamp				int					`json:"timestamp"`
//}

//type MozMmsMessage struct {
	//MozGenericMessage
	//subject						string
	//smil							interface{}
	//attachments				interface{}
	//expiryDate				int
//}

//func (m *MozMmsMessage) GetSubject() string {
	//return m.subject
//}

//func (m *MozMmsMessage) setSubject(subject string) string {
	//m.subject = subject
	//return m.subject
//}

//func (m *MozMmsMessage) GetSmil() interface{} {
	//return m.smil
//}

//func (m *MozMmsMessage) setSmil(smil interface{}) interface{} {
	//m.smil = smil
	//return m.smil
//}

//func (m *MozMmsMessage) GetAttachments() interface{} {
	//return m.attachments
//}

//func (m *MozMmsMessage) setAttachments(attachments interface{}) interface{} {
	//m.attachments = attachments
	//return m.attachments
//}

//func (m *MozMmsMessage) GetExpiryDate() int {
	//return m.expiryDate
//}

//func (m *MozMmsMessage) setExpiryDate(expiryDate int) int {
	//m.expiryDate = expiryDate
	//return m.expiryDate
//}

//func (m MozMmsMessage) MarshalJSON() ([]byte, error) {
	//j, err := json.Marshal(publicMozMmsMessage{
		//MozType: m.mozType,
		//Id: m.id,
		//ThreadId: m.threadId,
		//Subject: m.subject,
		//Smil: m.smil,
		//Attachments: m.attachments,
		//ExpiryDate: m.expiryDate,
		//Delivery: m.delivery,
		//DeliveryStatus: m.deliveryStatus,
		//Read: m.read,
		//Receiver: m.receiver,
		//Sender: m.sender,
		//Timestamp: m.timestamp,
	//})
	//if err != nil {
		//return nil, err
	//}
	//return j, nil
//}

//func (m MozMmsMessage) UnmarshalJSON(data []byte) (MozGenericMessageInterface, error) {
	//var cast *publicMozMmsMessage
	//if err := json.Unmarshal(data, &cast); err != nil {
		//return m, err
	//}
	//m.setType(cast.MozType)
	//m.setId(cast.Id)
	//m.setThreadId(cast.ThreadId)
	//m.setSubject(cast.Subject)
	//m.setSmil(cast.Smil)
	//m.setAttachments(cast.Attachments)
	//m.setExpiryDate(cast.ExpiryDate)
	//m.setDelivery(cast.Delivery)
	//m.setDeliveryStatus(cast.DeliveryStatus)
	//m.setRead(cast.Read)
	//m.setReceiver(cast.Receiver)
	//m.setSender(cast.Sender)
	//m.setTimestamp(cast.Timestamp)
	//return m, nil
//}

package types

type Class int

const (
	NORMAL	Class = iota
	CLASS_0
	CLASS_1
	CLASS_2
	CLASS_3
)

type MozSmsMessage struct {
	Type						string		`json:"type"`
	Id							int				`json:"id"`
	ThreadId				int				`json:"threadId"`
	IccId						string		`json:"iccId"`
	Body						string		`json:"body"`
	Delivery				string		`json:"delivery"`
	DeliveryStatus	string		`json:"deliveryStatus"`
	Read						bool			`json:"read"`
	Receiver				string		`json:"receiver"`
	Sender					string		`json:"sender"`
	Timestamp				int				`json:"timestamp"`
	MessageClass		string		`json:"class"`
}

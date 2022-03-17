package types

//import (
	//"encoding/json"
//)

type MozMobileMessageThread struct {
	Id									int				`json:"id"`
	Body								string		`json:"body"`
	UnreadCount					int				`json:"unreadCount"`
	Participants				[]string	`json:"participants"`
	Timestamp						int				`json:"timestamp"`
	LastMessageSubject	string		`json:"lastMessageSubject"`
	LastMessageType			string		`json:"lastMessageType"`
}

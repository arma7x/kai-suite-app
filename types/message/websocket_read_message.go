package message

import (
	"encoding/json"
)

type ReadKind int

// TODO
const (
	VERIFY_DEVICE 		ReadKind = iota
	DELETE_MESSAGE 		// id
	GET_MESSAGE 			// id
	GET_MESSAGES 			// filter, reverseOrder
	MARK_MESSAGE_READ // id, isRead
	SEND_SMS					// []number, message
	SEND_MMS					// param
	GET_THREADS
	RETRIEVE_MMS 			// id
	// CONTACT
	// CONTACT_SYNC // bulk
	// CALENDAR
	// CALENDAR_SYNC // bulk
	// CALL ??
)

type publicReadMessage struct {
	Action	string		`json:"action"`
	Kind 		ReadKind	`json:"kind"`
	Content	string		`json:"content"`
}

type ReadMessage struct {
	action 	string
	kind		ReadKind
	content string
}

func (r *ReadMessage) GetAction() string {
	return r.action
}

func (r *ReadMessage) setAction(action string) string {
	r.action = action
	return r.action
}

func (r *ReadMessage) GetKind() ReadKind {
	return r.kind
}

func (r *ReadMessage) setKind(kind ReadKind) ReadKind {
	r.kind = kind
	return r.kind
}

func (r *ReadMessage) GetContent() string {
	return r.content
}

func (r *ReadMessage) setContent(content string) string {
	r.content = content
	return r.content
}

func (r *ReadMessage) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicReadMessage{Action: r.action, Kind: r.kind, Content: r.content})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (r *ReadMessage) UnmarshalJSON(data []byte) (*ReadMessage, error) {
	var public *publicReadMessage
	if err := json.Unmarshal(data, &public); err != nil {
		return r, err
	}
	r.setAction(public.Action);
	r.setKind(public.Kind);
	r.setContent(public.Content);
	return r, nil
}

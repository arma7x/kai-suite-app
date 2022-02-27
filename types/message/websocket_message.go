package message

import (
	"encoding/json"
)

type Kind int64

// TODO
const (
	SMS						Kind = iota
	SMS_SYNC 			// bulk
	CONTACT
	CONTACT_SYNC	// bulk
	// CALENDAR
	// CALENDAR_SYNC // bulk
	// CALL ??
)

type publicWebsocketMessage struct {
	Action	string	`json:"action"`
	Kind 		Kind	`json:"kind"`
	Content	string	`json:"content"`
}

type WebsocketMessage struct {
	action 	string
	kind		Kind
	content string
}

func (r *WebsocketMessage) GetAction() string {
	return r.action
}

func (r *WebsocketMessage) setAction(action string) string {
	r.action = action
	return r.action
}

func (r *WebsocketMessage) GetKind() Kind {
	return r.kind
}

func (r *WebsocketMessage) setKind(kind Kind) Kind {
	r.kind = kind
	return r.kind
}

func (r *WebsocketMessage) GetContent() string {
	return r.content
}

func (r *WebsocketMessage) setContent(content string) string {
	r.content = content
	return r.content
}

func (r *WebsocketMessage) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicWebsocketMessage{Action: r.action, Kind: r.kind, Content: r.content})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (r *WebsocketMessage) UnmarshalJSON(data []byte) (*WebsocketMessage, error) {
	var public *publicWebsocketMessage
	if err := json.Unmarshal(data, &public); err != nil {
		return r, err
	}
	r.setAction(public.Action);
	r.setKind(public.Kind);
	r.setContent(public.Content);
	return r, nil
}

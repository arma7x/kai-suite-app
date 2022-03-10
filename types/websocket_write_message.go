package types

import (
	"encoding/json"
)

type WriteKind int

const (
	DELETE_MESSAGE		WriteKind = iota // content: id
	GET_MESSAGE 			// content: id
	GET_MESSAGES 			// content: filter, reverseOrder
	MARK_MESSAGE_READ // content: id, isRead
	SEND_SMS					// content: []number, message
	SEND_MMS					// content: param
	GET_THREADS
	RETRIEVE_MMS 			// content: id
)

type publicWriteMessage struct {
	Kind 		WriteKind	`json:"kind"`
	Content	string		`json:"content"`
}

type WriteMessage struct {
	kind		WriteKind
	content string
}

func (r *WriteMessage) GetKind() WriteKind {
	return r.kind
}

func (r *WriteMessage) setKind(kind WriteKind) WriteKind {
	r.kind = kind
	return r.kind
}

func (r *WriteMessage) GetContent() string {
	return r.content
}

func (r *WriteMessage) setContent(content string) string {
	r.content = content
	return r.content
}

func (r *WriteMessage) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicWriteMessage{ Kind: r.kind, Content: r.content})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (r *WriteMessage) UnmarshalJSON(data []byte) (*WriteMessage, error) {
	var public *publicWriteMessage
	if err := json.Unmarshal(data, &public); err != nil {
		return r, err
	}
	r.setKind(public.Kind);
	r.setContent(public.Content);
	return r, nil
}

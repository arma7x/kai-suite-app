package types

import (
	"encoding/json"
)

type ReadKind int

const (
	ACK				ReadKind = iota
	MESSAGE
	MESSAGES
	THREADS
	MMS
)

type publicReadMessage struct {
	Kind 		ReadKind	`json:"kind"`
	Content	string		`json:"content"`
}

type ReadMessage struct {
	kind		ReadKind
	content string
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
	j, err := json.Marshal(publicReadMessage{Kind: r.kind, Content: r.content})
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
	r.setKind(public.Kind);
	r.setContent(public.Content);
	return r, nil
}

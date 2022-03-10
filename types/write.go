package types

import (
	"encoding/json"
)

type publicWrite struct {
	Action 	string	`json:"action"`
	Content string	`json:"content"`
}

type Write struct {
	action string
	content string
}

func (w *Write) GetAction() string {
	return w.action
}

func (w *Write) SetAction(action string) string {
	w.action = action
	return w.action
}

func (w *Write) GetContent() string {
	return w.content
}

func (w *Write) SetContent(content string) string {
	w.content = content
	return w.content
}

func (w *Write) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicWrite{Action: w.action, Content: w.content})
	if err != nil {
		return nil, err
	}
	return j, nil
}

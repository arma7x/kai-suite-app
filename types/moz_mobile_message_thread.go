package types

import (
	"encoding/json"
)

type publicMozMobileMessageThread struct {
	Id							int				`json:"id"`
	Body						string		`json:"body"`
	UnreadCount			int				`json:"unreadCount"`
	Participants		[]string	`json:"participants"`
	Timestamp				int				`json:"timestamp"`
	LastMessageType	string		`json:"lastMessageType"`
}

type MozMobileMessageThread struct {
	id							int
	body						string
	unreadCount			int
	participants		[]string
	timestamp				int
	lastMessageType	string
}

func (m *MozMobileMessageThread) GetId() int {
	return m.id
}

func (m *MozMobileMessageThread) setId(id int) int {
	m.id = id
	return m.id
}

func (m *MozMobileMessageThread) GetBody() string {
	return m.body
}

func (m *MozMobileMessageThread) setBody(body string) string {
	m.body = body
	return m.body
}

func (m *MozMobileMessageThread) GetUnreadCount() int {
	return m.unreadCount
}

func (m *MozMobileMessageThread) setUnreadCount(unreadCount int) int {
	m.unreadCount = unreadCount
	return m.unreadCount
}

func (m *MozMobileMessageThread) GetParticipants() []string {
	return m.participants
}

func (m *MozMobileMessageThread) setParticipants(participants []string) []string {
	m.participants = participants
	return m.participants
}

func (m *MozMobileMessageThread) GetTimestamp() int {
	return m.timestamp
}

func (m *MozMobileMessageThread) setTimestamp(timestamp int) int {
	m.timestamp = timestamp
	return m.timestamp
}

func (m *MozMobileMessageThread) GetLastMessageType() string {
	return m.lastMessageType
}

func (m *MozMobileMessageThread) setLastMessageType(lastMessageType string) string {
	m.lastMessageType = lastMessageType
	return m.lastMessageType
}

func (m *MozMobileMessageThread) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(publicMozMobileMessageThread{
		Id: m.id,
		Body: m.body,
		UnreadCount: m.unreadCount,
		Participants: m.participants,
		Timestamp: m.timestamp,
		LastMessageType: m.lastMessageType,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (m *MozMobileMessageThread) UnmarshalJSON(data []byte) (*MozMobileMessageThread, error) {
	var public *publicMozMobileMessageThread
	if err := json.Unmarshal(data, &public); err != nil {
		return m, err
	}
	m.setId(public.Id);
	m.setBody(public.Body);
	m.setUnreadCount(public.UnreadCount);
	m.setParticipants(public.Participants);
	m.setTimestamp(public.Timestamp);
	m.setLastMessageType(public.LastMessageType);
	return m, nil
}

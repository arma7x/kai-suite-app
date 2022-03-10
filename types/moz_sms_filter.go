package types

import (
	"encoding/json"
)

type MozSmsFilter struct {
	StartDate	int				`json:"startDate"`
	EndDate		int				`json:"endDate"`
	Numbers		[]string	`json:"numbers"`
	Delivery	string		`json:"delivery"`
	Read			bool			`json:"read"`
	ThreadId	int				`json:"threadId"`
}

func (m *MozSmsFilter) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}

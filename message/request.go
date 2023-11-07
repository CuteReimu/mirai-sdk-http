package message

type SendGroupMessage struct {
	Target       int64 `json:"target,omitempty"`
	Group        int64 `json:"group,omitempty"`
	Quote        int64 `json:"quote,omitempty"`
	MessageChain Chain `json:"messageChain"`
}

func (m *SendGroupMessage) GetCommand() string {
	return "sendGroupMessage"
}

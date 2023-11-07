package message

type Chain []map[string]any

type SingleMessage interface {
	ToMap() Chain
}

type Plain struct {
	Text string
}

func (m *Plain) ToMap() map[string]any {
	return map[string]any{
		"type": "Plain",
		"text": m.Text,
	}
}

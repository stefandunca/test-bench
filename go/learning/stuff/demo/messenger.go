package demo

type Listener interface {
	Notify(string)
}

type Messenger struct {
	listeners []Listener
}

func (m *Messenger) Register(listener Listener) {
	m.listeners = append(m.listeners, listener)
}

func (m *Messenger) Send(message string) {
	for _, listener := range m.listeners {
		(listener).Notify(message)
	}
}

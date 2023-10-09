package server

type manager struct {
	lobs map[string]*lobby
}

func NewManager() *manager{
	m := manager{

	}
	return &m
}

func (m *manager) FindLobby(id string) *lobby{
	l, ok := m.lobs[id]
	if !ok {
		return nil
	}
	return l
}

func (m *manager) CreateLobby() *lobby{
	done := make(chan string, 1)
	go func() {
		id := <-done
		delete(m.lobs, id)
		close(done)
	}()

	l, id := NewLobby(done)
	m.lobs[id] = l
	return l
}
package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math"
)

type manager struct {
	lobs map[string]*lobby
}

func NewManager() *manager{
	m := manager{
		lobs: make(map[string]*lobby),
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

	id := m.randId(6)
	l := NewLobby(done, id)
	m.lobs[id] = l
	return l
}

func (m *manager) randId(length int) string {
	// https://stackoverflow.com/a/75518426/7361588
	for {
		bi, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(math.Pow(10, float64(length)))),
		)
		if err != nil {
			panic(err)
		}
		id := fmt.Sprintf("%0*d", length, bi)
		_, ok := m.lobs[id]
		if !ok {
			return id
		}
	}
}


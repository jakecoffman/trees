package lib

import (
	"github.com/gorilla/websocket"
	"sync"
)

type SafetySocket struct {
	rxLock sync.Mutex
	txLock sync.Mutex
	ws     *websocket.Conn
}

func NewSafetySocket(ws *websocket.Conn) *SafetySocket {
	return &SafetySocket{ws: ws}
}

func (s *SafetySocket) ReadJSON(v interface{}) error {
	s.rxLock.Lock()
	defer s.rxLock.Unlock()

	return s.ws.ReadJSON(v)
}

func (s *SafetySocket) WriteJSON(v interface{}) error {
	s.txLock.Lock()
	defer s.txLock.Unlock()

	return s.ws.WriteJSON(v)
}

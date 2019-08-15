package websocket

import (
	"io"
	"time"

	"github.com/gorilla/websocket"
)

type wsconn struct {
	*websocket.Conn
	reader io.Reader
}

func (wc *wsconn) Read(data []byte) (int, error) {
	if wc.reader == nil {
		_, reader, err := wc.NextReader()
		if err != nil {
			return 0, err
		}

		wc.reader = reader
	}

	n, err := wc.reader.Read(data)
	if err != nil {
		wc.reader = nil
	}

	return n, err
}

func (wc *wsconn) Write(data []byte) (int, error) {
	err := wc.WriteMessage(websocket.BinaryMessage, data)
	return 0, err
}

func (wc *wsconn) SetDeadline(t time.Time) error {
	if err := wc.SetReadDeadline(t); err != nil {
		return err
	}

	return wc.SetWriteDeadline(t)
}

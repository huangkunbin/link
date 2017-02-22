package link

import (
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

type NetConn struct {
	ws *websocket.Conn
	r  io.Reader
}

func (conn *NetConn) Read(b []byte) (n int, err error) {
	var opCode int
	if conn.r == nil {
		var r io.Reader
		for {
			if opCode, r, err = conn.ws.NextReader(); err != nil {
				return
			}
			if opCode != websocket.TextMessage && opCode != websocket.BinaryMessage {
				continue
			}

			conn.r = r
			break
		}
	}

	n, err = conn.r.Read(b)
	if err != nil {
		if err == io.EOF {
			conn.r = nil
			err = nil
		}
	}
	return
}

func (conn *NetConn) Write(b []byte) (n int, err error) {
	var w io.WriteCloser
	if w, err = conn.ws.NextWriter(websocket.TextMessage); err != nil {
		return
	}
	if n, err = w.Write(b); err != nil {
		return
	}
	err = w.Close()
	return
}

func (conn *NetConn) Close() error {
	return conn.ws.Close()
}

func (conn *NetConn) LocalAddr() net.Addr {
	return conn.ws.LocalAddr()
}

func (conn *NetConn) RemoteAddr() net.Addr {
	return conn.ws.RemoteAddr()
}

func (conn *NetConn) SetDeadline(t time.Time) (err error) {
	if err = conn.ws.SetReadDeadline(t); err != nil {
		return
	}
	return conn.ws.SetWriteDeadline(t)
}

func (conn *NetConn) SetReadDeadline(t time.Time) error {
	return conn.ws.SetReadDeadline(t)
}

func (conn *NetConn) SetWriteDeadline(t time.Time) error {
	return conn.ws.SetWriteDeadline(t)
}

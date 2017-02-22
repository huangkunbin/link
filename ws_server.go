package link

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSHandler struct {
	upgrader websocket.Upgrader
	server   *Server
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade error: %v", err)
		return
	}

	netConn := &NetConn{
		ws: conn,
	}

	go func() {
		codec, err := handler.server.protocol.NewCodec(netConn)
		if err != nil {
			netConn.Close()
			return
		}

		log.Println("new session")
		session := handler.server.manager.NewSession(codec, handler.server.sendChanSize)
		handler.server.handler.HandleSession(session)
	}()

}

func (server *Server) WSServe() {
	httpServer := &http.Server{
		Addr: server.Listener().Addr().String(),
		Handler: &WSHandler{upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
			server: server,
		},
	}

	httpServer.Serve(server.listener)
}

func WSDial(address string, protocol Protocol, sendChanSize int) (*Session, error) {
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}
	netConn := &NetConn{
		ws: conn,
	}
	codec, err := protocol.NewCodec(netConn)
	if err != nil {
		return nil, err
	}
	return NewSession(codec, sendChanSize), nil
}

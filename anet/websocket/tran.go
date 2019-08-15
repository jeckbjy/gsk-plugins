package websocket

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeckbjy/gsk/anet"
	"github.com/jeckbjy/gsk/anet/base"
)

func init() {
	anet.Add("websocket", New)
}

func New() anet.ITran {
	return &wstran{}
}

type wstran struct {
	base.Tran
	websocket.Upgrader
}

func (*wstran) String() string {
	return "websocket"
}

func (t *wstran) Listen(addr string, opts ...anet.ListenOption) (anet.IListener, error) {
	conf := anet.ListenOptions{}
	conf.Init(opts...)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := http.Server{Handler: &handler{upgrader: &t.Upgrader, tran: t, tag: conf.Tag}}
	go func() {
		_ = server.Serve(l)
	}()
	return l, nil
}

func (t *wstran) Dial(addr string, opts ...anet.DialOption) (anet.IConn, error) {
	conf := &anet.DialOptions{}
	conf.Init(opts...)
	return base.DoDial(conf, t, func() (net.Conn, error) {
		//
		ctx := context.Background()
		if conf.Timeout != 0 {
			ctx, _ = context.WithTimeout(ctx, conf.Timeout)
		}

		c, _, err := websocket.DefaultDialer.DialContext(ctx, addr, nil)
		if err != nil {
			return nil, err
		}

		return &wsconn{Conn: c}, nil
	})
}

type handler struct {
	upgrader *websocket.Upgrader
	tran     *wstran
	tag      string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err == nil {
		c := &wsconn{Conn: conn}
		base.DoOpen(c, h.tran, false, h.tag)
	}
}

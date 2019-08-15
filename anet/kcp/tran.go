package kcp

import (
	"net"

	"github.com/jeckbjy/gsk/anet"
	"github.com/jeckbjy/gsk/anet/base"
	kcpgo "github.com/xtaci/kcp-go"
)

func init() {
	anet.Add("kcp", New)
}

func New() anet.ITran {
	return &ktran{}
}

type ktran struct {
	base.Tran
}

func (t *ktran) String() string {
	return "kcp"
}

func (t *ktran) Listen(addr string, opts ...anet.ListenOption) (anet.IListener, error) {
	conf := anet.ListenOptions{}
	conf.Init(opts...)
	// TODO:add kcp option
	return base.DoListen(&conf, t, func() (net.Listener, error) {
		return kcpgo.Listen(addr)
	})
}

func (t *ktran) Dial(addr string, opts ...anet.DialOption) (anet.IConn, error) {
	conf := &anet.DialOptions{}
	conf.Init(opts...)
	return base.DoDial(conf, t, func() (net.Conn, error) {
		// TODO:timeout,options
		return kcpgo.Dial(addr)
	})
}

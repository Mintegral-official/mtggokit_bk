package bifrost

import (
	"github.com/Mintegral-official/mtggokit/bifrost/log"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
)

type BifrostServerCfg struct {
	Addr   string
	Logger log.BiLogger
}

type BifrostServer struct {
	cfg            *BifrostServerCfg
	server         *server.Server
	ProviderManger *streamer.StreamerProviderManager
}

func NewBifrostServer(cfg *BifrostServerCfg) *BifrostServer {
	s := server.NewServer()
	spm := streamer.NewStreamerProviderManager()
	biServer := streamer.NewBifrostServer(spm)
	share.Codecs[protocol.SerializeType(5)] = &streamer.GobCodec{}
	if e := s.RegisterName("BifrostService", biServer, ""); e != nil {
		cfg.Logger.Warn("RegisterName BifrostService error")
		return nil
	}
	return &BifrostServer{
		cfg:            cfg,
		server:         s,
		ProviderManger: spm,
	}
}

func (bs *BifrostServer) Serve() error {
	return bs.server.Serve("tcp", bs.cfg.Addr)
}

func (bs *BifrostServer) Close() error {
	return bs.server.Close()
}

func (bs *BifrostServer) RegisterProvider(name string, sp *streamer.StreamerProvider) error {
	if e := bs.ProviderManger.RegiterProvider(name, sp); e != nil {
		if bs.cfg.Logger != nil {
			bs.cfg.Logger.Warn("Register error:" + e.Error())
		}
		return e
	}
	return nil
}

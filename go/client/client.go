package client

import (
	"context"
	"github.com/obnahsgnaw/socketclient/go/base"
	"github.com/obnahsgnaw/socketutil/service/client"
	"go.uber.org/zap/zapcore"
	"strconv"
)

type Client struct {
	base.Server
	ctx    context.Context
	client *client.Client
	config *Config
}

func New(ctx context.Context, cc *Config) *Client {
	cc.parseAddr()
	c := client.New(ctx, cc.network, cc.ip+":"+strconv.Itoa(cc.port)+cc.path, cc.ProtocolCoder, cc.GatewayPkgCoder, cc.DataCoder)
	s := &Client{
		ctx:    ctx,
		client: c,
		config: cc,
	}
	c.With(client.Connect(func(index int) {
		if s.config.ServerLogWatcher != nil {
			s.config.ServerLogWatcher(zapcore.InfoLevel, "client connected")
		}
		s.Ready()
	}))
	c.With(client.Disconnect(func(index int) {
		if s.config.ServerLogWatcher != nil {
			s.config.ServerLogWatcher(zapcore.InfoLevel, "client disconnected")
		}
		s.Pause()
	}))
	if cc.RetryInterval > 0 {
		c.With(client.Retry(cc.RetryInterval))
	}
	if cc.KeepaliveInterval > 0 {
		c.With(client.Keepalive(cc.KeepaliveInterval))
	}
	if cc.Timeout > 0 {
		c.With(client.Timeout(cc.Timeout))
	}
	if cc.ServerLogWatcher != nil {
		c.With(client.Logger(cc.ServerLogWatcher))
	}
	if cc.ActionLogWatcher != nil {
		c.With(client.ActionLogger(cc.ActionLogWatcher))
	}
	if cc.PackageLogWatcher != nil {
		c.With(client.PackageLogger(cc.PackageLogWatcher))
	}

	return s
}

func (s *Client) Start() {
	s.client.Start()
}

func (s *Client) Stop() {
	s.client.Stop()
	s.Pause()
}

func (s *Client) Client() *client.Client {
	return s.client
}

func (s *Client) Context() context.Context {
	return s.ctx
}

func (s *Client) Config() *Config {
	return s.config
}

func (s *Client) Log(l zapcore.Level, msg string) {
	if s.config.ServerLogWatcher != nil {
		s.config.ServerLogWatcher(l, msg)
	}
}

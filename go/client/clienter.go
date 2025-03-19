package client

import (
	"context"
	"github.com/obnahsgnaw/socketutil/service/client"
	"go.uber.org/zap/zapcore"
)

type Clienter interface {
	Start()
	Stop()
	Client() *client.Client
	Context() context.Context
	Config() *Config
	Log(l zapcore.Level, msg string)
	WhenReady(cb func())
	WhenPaused(cb func())
}

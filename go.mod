module github.com/obnahsgnaw/socketclient

go 1.19

require (
	github.com/obnahsgnaw/application v0.17.16
	github.com/obnahsgnaw/socketutil v0.8.9
	go.uber.org/zap v1.23.0
	google.golang.org/protobuf v1.33.0
)

replace github.com/obnahsgnaw/socketutil v0.8.8 => ../socketutil

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

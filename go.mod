module github.com/fgrosse/go-home

go 1.12

require (
	github.com/faiface/pixel v0.8.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.4
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/image v0.0.0-20190523035834-f03afa92d3ff
	gopkg.in/yaml.v3 v3.0.0-20190502103701-55513cacd4ae
)

replace github.com/faiface/pixel => github.com/fgrosse/pixel v0.8.1-0.20190530122943-dd5e0d8b09c3

package init

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/example/app/foobar"
)

func init() {
	err := com.DefaultRegistry.Register(
		&com.Object{Value: &foobar.Component{}})
	if err != nil {
		panic(err)
	}
}

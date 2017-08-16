package main

import (
	"github.com/gliderlabs/com/daemon"

	_ "github.com/gliderlabs/com/example/app/baz"
	_ "github.com/gliderlabs/com/example/app/foobar/init"
	_ "github.com/gliderlabs/com/example/app/qux"
)

func main() {
	daemon.Run("example")
}

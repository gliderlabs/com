package main

import "github.com/gliderlabs/com/example/app/ssh"

func main() {}

func Registerable() []interface{} {
	return []interface{}{
		&ssh.Component{},
	}
}

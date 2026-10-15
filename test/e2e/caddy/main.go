package main

import (
	"github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/danmoz/caddy-lambda-licious"
)

func main() {
	caddycmd.Main()
}

package main

import (
	"context"
	"goexiftool/pkg/providers"
	"goexiftool/pkg/servers"
)

const (
	defaultPort = 8080
)

//create webserver and listen on port 8080 using standard http
func main() {
	provider := providers.New()
	server := servers.New()
	if err := server.Run(context.Background(), provider, defaultPort); err != nil {
		panic(err)
	}
}

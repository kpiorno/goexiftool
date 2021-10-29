package pkg

import "context"

//Server defines the interface for the webserver
type Server interface {
	Run(ctx context.Context, provider Provider, httpPort int) error
}

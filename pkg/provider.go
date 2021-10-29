package pkg

import (
	"goexiftool/pkg/dtos"
)

//Provider defines the interface for the provider
type Provider interface {
	ConsumeFile(file string) (<-chan *dtos.Response, error)
	Close()
}

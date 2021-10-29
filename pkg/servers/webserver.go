package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"goexiftool/pkg"
	"goexiftool/pkg/dtos"
	"net/http"
)

var _ pkg.Server = &WebServer{}

type WebServer struct {
	provider        pkg.Provider
	currentResponse *dtos.Response
	command         chan Command
}

const (
	registerConsumer int = iota
)

type Command struct {
	CommandType int
	Response    chan *dtos.Response
}

func (s *WebServer) handler(w http.ResponseWriter, r *http.Request) {
	data := make(chan *dtos.Response)
	s.command <- Command{CommandType: registerConsumer, Response: data}
	response := <-data
	jsonData, _ := json.Marshal(response)
	if _, err := w.Write(jsonData); err != nil {
		panic(err)
	}

}

func (s *WebServer) Run(ctx context.Context, provider pkg.Provider, httpPort int) error {
	s.provider = provider
	providerChan, err := provider.ConsumeFile("-listx")
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				provider.Close()
				return
			case providerData := <-providerChan:
				s.currentResponse = providerData
			case command := <-s.command:
				switch command.CommandType {
				case registerConsumer:
					command.Response <- s.currentResponse
				}
			}
		}
	}()
	http.HandleFunc("/tags", s.handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		return err
	}
	return nil
}

func New() *WebServer {
	return &WebServer{
		command: make(chan Command),
	}
}

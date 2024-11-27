package app

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
}

func (s *GreetService) OnStartup(ctx application.Context, options application.ServiceOptions) error {
	return nil
}

func (s *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}

func (s *GreetService) SyncCheck(timestamp string) {
	log.Default().Printf("SyncCheck: %s\n", timestamp)
}

package main

import (
	"github.com/kardianos/service"
)

type daemon struct{}

func (p *daemon) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *daemon) run() {
	run()
}

func (p *daemon) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "podcastd",
		DisplayName: "podcastd",
		Description: "https://github.com/bigwhoop/podcastd (v" + VERSION + ")",
	}

	prg := &daemon{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}

	if s.Run() != nil {
		logger.Fatal(err)
	}
}

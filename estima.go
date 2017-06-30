package main

import (
	"ru/sbt/estima/app"
	"github.com/kardianos/service"
	"os"
	"log"
	"fmt"
)

type Program struct {}

func (p *Program) Start (s service.Service) error {
	go p.run()
	return nil
}

func (p *Program) run () {
	app.AppRun()
}

func (p *Program) Stop (s service.Service) error {
	go p.run()
	return nil
}

func main() {
	if len(os.Args) == 0 {
		app.AppRun()
	} else {
		serviceName := os.Args[1]
		log.Printf("Install service with name: %v\n", serviceName)
		svcConfig := &service.Config{
			Name: serviceName,
			DisplayName: "Estimator service",
			Description: "Estimator service",
			Dependencies: []string{"ArangoDB"},
		}

		prg := &Program{}
		srv, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}

		logger, err := srv.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = srv.Run()
		if err != nil {
			logger.Error(err)
		}
	}
}

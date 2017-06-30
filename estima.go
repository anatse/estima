package main

import (
	"ru/sbt/estima/app"
	"github.com/kardianos/service"
	"os"
	"log"
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
	return nil
}

func main() {
	f, err := os.OpenFile("estima.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting estima")

	if len(os.Args) == 1 {
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

		err = srv.Install()
		if err != nil {
			logger.Error(err)
		}

		err = srv.Start()
		if err != nil {
			logger.Error(err)
		}

		log.Printf("Installed service with name: %v\n", serviceName)
	}
}

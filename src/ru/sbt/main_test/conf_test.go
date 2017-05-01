package main_test

import (
	"testing"
	"ru/sbt/estima/conf"
	"os"
)

func TestConfog(t *testing.T) {
	config := conf.LoadConfig()
	if config.Name == "" {
		t.Errorf("Configuration not loaded. Environment: %s", os.Getenv("CONFIG_PATH"))
	}

	if config.Secret != "secret" {
		t.Error("Dont change secret key in development mode!")
	}

	if config.Database.Name != "estima" {
		t.Error("Database name should be estima")
	}

	if config.Database.Password == "" {
		t.Error("Password not defined")
	}
}
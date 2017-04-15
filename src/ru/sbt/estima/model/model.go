package model

import (
	ara "github.com/diegogub/aranGO"
)

type Entity interface {
	AraDoc() (ara.Document)
	Entity() interface{}

	FromJson(json []byte)(Entity, error)
	Copy (Entity)Entity
}


package model

import (
	ara "github.com/diegogub/aranGO"
)

type Entity interface {
	AraDoc() (ara.Document)
	Entity() interface{}
	GetKey() string
	GetCollection() string
}

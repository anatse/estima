package model

import (
	ara "github.com/diegogub/aranGO"
)

type Entity interface {
	AraDoc() (ara.Document)
	ToJson() ([]byte, error)
	FromJson(json []byte)(error)
	Copy (Entity)
}


package model

import (
	ara "github.com/diegogub/aranGO"
	"time"
)

type CalcUnit struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	Version string `json:"version,omitempty"`
	Status string `json:"status,omitempty"`
	Changed time.Time `json:"changed,omitempty"`
}

func (cmp CalcUnit) Entity() interface{} {
	return struct{
		*CalcUnit

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&cmp,
		nil,
		nil,
		nil,
		nil,
	}
}

func (cmp CalcUnit) AraDoc() (ara.Document) {
	return cmp.Document
}

func (cmp CalcUnit)GetKey() string {
	return cmp.Key
}

func (cmp CalcUnit) GetCollection() string {
	return "calcunits"
}

func (cmp CalcUnit) GetError()(string, bool) {
	return cmp.Message, cmp.Error
}
package model

import (
	ara "github.com/diegogub/aranGO"
	"math/big"
	"time"
)

type CalcUnitPrice struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	StoryPoints big.Float `json:"storyPoins,omitempty"`
	Group string `json:"group,omitempty"`
	Changed time.Time `json:"changed,omitempty"`
}

func (cmp CalcUnitPrice) Entity() interface{} {
	return struct{
		*CalcUnitPrice

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

func (cmp CalcUnitPrice) AraDoc() (ara.Document) {
	return cmp.Document
}

func (cmp CalcUnitPrice)GetKey() string {
	return cmp.Key
}

func (cmp CalcUnitPrice) GetCollection() string {
	return "cuprices"
}

func (cmp CalcUnitPrice) GetError()(string, bool) {
	return cmp.Message, cmp.Error
}
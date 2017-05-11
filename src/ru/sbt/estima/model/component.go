package model

import (
	ara "github.com/diegogub/aranGO"
)

type Component struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty,required"`
	Description string `json:"description,omitempty"`
	Owner string `json:"owner,omitempty"`
}

func (cmp Component) Entity() interface{} {
	return struct{
		*Component

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

func (cmp Component) AraDoc() (ara.Document) {
	return cmp.Document
}

func (cmp Component)GetKey() string {
	return cmp.Key
}

func (cmp Component) GetCollection() string {
	return "components"
}

func (cmp Component) GetError()(string, bool) {
	return cmp.Message, cmp.Error
}
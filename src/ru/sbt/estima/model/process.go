package model

import (
	ara "github.com/diegogub/aranGO"
)

type Process struct {
	ara.Document
	Name string
	Description string
	Status string
}

func (prc Process) Entity() interface{} {
	return struct{
		*Process

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`
		OmitKey omit `json:"_key,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&prc,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (prc Process) AraDoc() (ara.Document) {
	return prc.Document
}

func (prc Process)GetKey() string {
	return prc.Name
}

func (prc Process) GetCollection() string {
	return "processes"
}

func (prc Process) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return prc.Message, prc.Error
}
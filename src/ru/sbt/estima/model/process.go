package model

import (
	ara "github.com/diegogub/aranGO"
)

type Process struct {
	ara.Document
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status string `json:"status,omitempty"`
}

func (prc Process) Entity() interface{} {
	return struct{
		*Process

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`


		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&prc,
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
	return prc.Key
}

func (prc Process) GetCollection() string {
	return "processes"
}

func (prc Process) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return prc.Message, prc.Error
}

func (prc Process) CopyChanged (entity Entity) Entity {
	newPrc := entity.(Process)
	if newPrc.Name != "" {prc.Name = newPrc.Name}
	if newPrc.Description != "" {prc.Description = newPrc.Description}
	if newPrc.Status != "" {prc.Status = newPrc.Status}
	return prc
}
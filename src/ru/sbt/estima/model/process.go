package model

import (
	ara "github.com/diegogub/aranGO"
	"encoding/json"
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

func (prc Process) Copy (entity Entity) Entity {
	var from Process = entity.(Process)
	prc.Name = from.Name
	prc.Description = from.Description
	return prc
}

func (prc Process) FromJson (jsUser []byte) (Entity, error) {
	var retprc Process
	err := json.Unmarshal(jsUser, &retprc)
	if err != nil {
		panic(err)
	}

	return retprc, err
}

func (prc Process)GetKey() string {
	return prc.Name
}

func (prc Process) GetCollection() string {
	return "Processs"
}

func (prc Process) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return prc.Message, prc.Error
}
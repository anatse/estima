package model

import (
	ara "github.com/diegogub/aranGO"
	"encoding/json"
)

type Feature struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	User string `json:"user,omitempty"`

}

func (fea Feature) Entity() interface{} {
	return struct{
		*Feature

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`
		OmitKey omit `json:"_key,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&fea,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (fea Feature) AraDoc() (ara.Document) {
	return fea.Document
}

func (fea Feature) Copy (entity Entity) Entity {
	var from Process = entity.(Process)
	fea.Name = from.Name
	return fea
}

func (fea Feature) FromJson (jsUser []byte) (Entity, error) {
	var retprc Process
	err := json.Unmarshal(jsUser, &retprc)
	if err != nil {
		panic(err)
	}

	return retprc, err
}

func (fea Feature)GetKey() string {
	return fea.Name
}

func (fea Feature) GetCollection() string {
	return "Processs"
}

func (fea Feature) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return fea.Message, fea.Error
}
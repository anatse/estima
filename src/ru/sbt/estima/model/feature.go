package model

import (
	ara "github.com/diegogub/aranGO"
)

type Feature struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty" required`
	Text string `json:"text,omitempty"`
	User string `json:"user,omitempty"`
	Version int `json:"version" required`
}

func (fea Feature) Entity() interface{} {
	return struct{
		*Feature

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&fea,
		nil,
		nil,
		nil,
		nil,
	}
}

func (fea Feature) AraDoc() (ara.Document) {
	return fea.Document
}

func (fea Feature)GetKey() string {
	return fea.Key
}

func (fea Feature) GetCollection() string {
	return "features"
}

func (fea Feature) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return fea.Message, fea.Error
}
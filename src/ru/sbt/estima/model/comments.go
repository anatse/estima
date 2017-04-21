package model

import (
	ara "github.com/diegogub/aranGO"
)

type Comment struct {
	ara.Document `json:-`
	Name string `json:"name,omitempty" required`
	Text string `json:"text,omitempty"`
	User string `json:"user,omitempty"`
	Version int `json:"version" required`
}

func (com Comment) Entity() interface{} {
	return struct{
		*Comment

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`
		OmitKey omit `json:"_key,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&com,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

func (com Comment) AraDoc() (ara.Document) {
	return com.Document
}

func (com Comment)GetKey() string {
	return com.Name
}

func (com Comment) GetCollection() string {
	return "comments"
}

func (com Comment) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return com.Message, com.Error
}

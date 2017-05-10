package model

import (
	ara "github.com/diegogub/aranGO"
)

type Comment struct {
	ara.Document `json:-`
	Title string `json:"name,omitempty"`
	Text string `json:"text,omitempty,required"`
}

func (com Comment) Entity() interface{} {
	return struct{
		*Comment

		OmitId  omit `json:"_id,omitempty"`
		OmitRev omit `json:"_rev,omitempty"`

		OmitError   omit   `json:"error,omitempty"`
		OmitMessage omit `json:"errorMessage,omitempty"`
	} {
		&com,
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
	return com.Key
}

func (com Comment) GetCollection() string {
	return "comments"
}

func (com Comment) GetError()(string, bool) {
	// default error bool and messages. Could be any kind of error
	return com.Message, com.Error
}
